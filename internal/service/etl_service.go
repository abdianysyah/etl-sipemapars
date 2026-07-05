package service

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	source "sipemapars-etl/internal/models/source"
	// warehouse "sipemapars-etl/internal/models/warehouse"
	"sipemapars-etl/internal/reporter"
	"sipemapars-etl/internal/repository"
	"sipemapars-etl/internal/transformer"
	"sipemapars-etl/internal/utils"
)

type Service struct {
	sourceRepo    *repository.SourceRepository
	warehouseRepo *repository.WarehouseRepository
	reporter      reporter.Reporter
}

func New(
	sourceRepo *repository.SourceRepository,
	warehouseRepo *repository.WarehouseRepository,
	rep reporter.Reporter,
) *Service {
	if rep == nil {
		rep = reporter.Noop{}
	}

	return &Service{
		sourceRepo:    sourceRepo,
		warehouseRepo: warehouseRepo,
		reporter:      rep,
	}
}

type patientMeta struct {
	DWID         int64
	SourceDietID int64
	DWRoomID     int64
	DWDietID     int64
}

func (s *Service) Run(ctx context.Context, jobUUID string) error {
	fmt.Println("etl start")

	_ = s.reporter.Log(jobUUID, 0, "Start", "ETL started")
	_ = s.reporter.Progress(jobUUID, 0, "Starting ETL")

	// =========================
	// EXTRACT
	// =========================

	_ = s.reporter.Progress(jobUUID, 5, "Extract Diets")
	diets, err := s.sourceRepo.GetDiets(ctx)
	if err != nil {
		_ = s.reporter.Failed(jobUUID, "Extract Diets", err.Error())
		return fmt.Errorf("extract diets: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Extract Diets", fmt.Sprintf("%d diets extracted", len(diets)))

	_ = s.reporter.Progress(jobUUID, 10, "Extract Rooms")
	rooms, err := s.sourceRepo.GetRooms(ctx)
	if err != nil {
		_ = s.reporter.Failed(jobUUID, "Extract Rooms", err.Error())
		return fmt.Errorf("extract rooms: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Extract Rooms", fmt.Sprintf("%d rooms extracted", len(rooms)))

	_ = s.reporter.Progress(jobUUID, 15, "Extract Users")
	users, err := s.sourceRepo.GetUsers(ctx)
	if err != nil {
		_ = s.reporter.Failed(jobUUID, "Extract Users", err.Error())
		return fmt.Errorf("extract users: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Extract Users", fmt.Sprintf("%d users extracted", len(users)))

	_ = s.reporter.Progress(jobUUID, 20, "Extract Menus")
	menus, err := s.sourceRepo.GetMenus(ctx)
	if err != nil {
		_ = s.reporter.Failed(jobUUID, "Extract Menus", err.Error())
		return fmt.Errorf("extract menus: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Extract Menus", fmt.Sprintf("%d menus extracted", len(menus)))

	_ = s.reporter.Progress(jobUUID, 30, "Extract Patients")
	patients, err := s.sourceRepo.GetPatients(ctx)
	if err != nil {
		_ = s.reporter.Failed(jobUUID, "Extract Patients", err.Error())
		return fmt.Errorf("extract patients: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Extract Patients", fmt.Sprintf("%d patients extracted", len(patients)))

	_ = s.reporter.Progress(jobUUID, 45, "Extract Transactions")
	transactions, err := s.sourceRepo.GetTransactions(ctx)
	if err != nil {
		_ = s.reporter.Failed(jobUUID, "Extract Transactions", err.Error())
		return fmt.Errorf("extract transactions: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Extract Transactions", fmt.Sprintf("%d transactions extracted", len(transactions)))

	// =========================
	// RESET WAREHOUSE
	// =========================

	_ = s.reporter.Progress(jobUUID, 55, "Reset Warehouse")
	if err := s.warehouseRepo.Reset(ctx); err != nil {
		_ = s.reporter.Failed(jobUUID, "Reset Warehouse", err.Error())
		return fmt.Errorf("reset warehouse: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Reset Warehouse", "Warehouse truncated successfully")

	// =========================
	// BEGIN TRANSACTION
	// =========================

	tx, err := s.warehouseRepo.BeginTx(ctx)
	if err != nil {
		_ = s.reporter.Failed(jobUUID, "Begin Transaction", err.Error())
		return fmt.Errorf("begin warehouse tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// =========================
	// LOAD DIMENSIONS
	// =========================

	_ = s.reporter.Progress(jobUUID, 60, "Load Dimensions")

	dietMap := make(map[int64]int64)
	roomMap := make(map[int64]int64)
	userMap := make(map[int64]int64)
	menuMapBySourceDiet := make(map[int64]int64)
	patientMetaBySourceID := make(map[string]patientMeta)
	timeMap := make(map[string]int64)

	if err := s.loadDiets(ctx, tx, diets, dietMap); err != nil {
		_ = s.reporter.Failed(jobUUID, "Load Diets", err.Error())
		return fmt.Errorf("load diets: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Load Diets", fmt.Sprintf("%d diets loaded", len(diets)))

	if err := s.loadRooms(ctx, tx, rooms, roomMap); err != nil {
		_ = s.reporter.Failed(jobUUID, "Load Rooms", err.Error())
		return fmt.Errorf("load rooms: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Load Rooms", fmt.Sprintf("%d rooms loaded", len(rooms)))

	if err := s.loadUsers(ctx, tx, users, userMap); err != nil {
		_ = s.reporter.Failed(jobUUID, "Load Users", err.Error())
		return fmt.Errorf("load users: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Load Users", fmt.Sprintf("%d users loaded", len(users)))

	if err := s.loadMenus(ctx, tx, menus, menuMapBySourceDiet); err != nil {
		_ = s.reporter.Failed(jobUUID, "Load Menus", err.Error())
		return fmt.Errorf("load menus: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Load Menus", fmt.Sprintf("%d menus loaded", len(menus)))

	if err := s.loadPatients(ctx, tx, patients, roomMap, dietMap, patientMetaBySourceID); err != nil {
		_ = s.reporter.Failed(jobUUID, "Load Patients", err.Error())
		return fmt.Errorf("load patients: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Load Patients", fmt.Sprintf("%d patients loaded", len(patients)))

	if err := s.loadTimes(ctx, tx, transactions, timeMap); err != nil {
		_ = s.reporter.Failed(jobUUID, "Load Time Dimension", err.Error())
		return fmt.Errorf("load time dimension: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Load Time Dimension", fmt.Sprintf("%d time rows loaded", len(timeMap)))

	_ = s.reporter.Progress(jobUUID, 85, "Load Fact Transactions")
	if err := s.loadFacts(ctx, tx, transactions, patientMetaBySourceID, menuMapBySourceDiet, timeMap); err != nil {
		_ = s.reporter.Failed(jobUUID, "Load Fact Transactions", err.Error())
		return fmt.Errorf("load facts: %w", err)
	}
	_ = s.reporter.Log(jobUUID, 0, "Load Fact Transactions", fmt.Sprintf("%d fact rows loaded", len(transactions)))

	// =========================
	// COMMIT
	// =========================

	if err := tx.Commit(); err != nil {
		_ = s.reporter.Failed(jobUUID, "Commit", err.Error())
		return fmt.Errorf("commit warehouse tx: %w", err)
	}

	_ = s.reporter.Progress(jobUUID, 100, "Completed")
	_ = s.reporter.Log(jobUUID, 3, "Completed", "ETL completed successfully")
	_ = s.reporter.Finish(jobUUID, "ETL completed successfully")

	fmt.Println("etl done")
	return nil
}

func (s *Service) loadDiets(ctx context.Context, tx *sql.Tx, rows []source.Diet, dietMap map[int64]int64) error {
	nextID := int64(2)

	for _, row := range rows {
		if err := s.warehouseRepo.InsertDiet(ctx, tx, transformer.Diet(nextID, row)); err != nil {
			return err
		}
		dietMap[row.ID] = nextID
		nextID++
	}

	return nil
}

func (s *Service) loadRooms(ctx context.Context, tx *sql.Tx, rows []source.Room, roomMap map[int64]int64) error {
	nextID := int64(2)

	for _, row := range rows {
		if err := s.warehouseRepo.InsertRoom(ctx, tx, transformer.Room(nextID, row)); err != nil {
			return err
		}
		roomMap[row.ID] = nextID
		nextID++
	}

	return nil
}

func (s *Service) loadUsers(ctx context.Context, tx *sql.Tx, rows []source.User, userMap map[int64]int64) error {
	nextID := int64(2)

	for _, row := range rows {
		if err := s.warehouseRepo.InsertUser(ctx, tx, transformer.User(nextID, row)); err != nil {
			return err
		}
		userMap[row.ID] = nextID
		nextID++
	}

	return nil
}

func (s *Service) loadMenus(ctx context.Context, tx *sql.Tx, rows []source.Menu, menuMapBySourceDiet map[int64]int64) error {
	nextID := int64(2)

	for _, row := range rows {
		if err := s.warehouseRepo.InsertMenu(ctx, tx, transformer.Menu(nextID, row)); err != nil {
			return err
		}

		// Simpan 1 mapping menu per diet asal.
		// Kalau 1 diet punya 1 menu, ini aman.
		if _, exists := menuMapBySourceDiet[row.DietID]; !exists {
			menuMapBySourceDiet[row.DietID] = nextID
		}

		nextID++
	}

	return nil
}

func (s *Service) loadPatients(
	ctx context.Context,
	tx *sql.Tx,
	rows []source.Patient,
	roomMap, dietMap map[int64]int64,
	patientMetaBySourceID map[string]patientMeta,
) error {
	nextID := int64(2)

	for _, row := range rows {
		roomID, ok := roomMap[row.RoomID]
		if !ok {
			return fmt.Errorf("room mapping not found for source room id=%d (patient=%s)", row.RoomID, row.ID)
		}

		dietID, ok := dietMap[row.DietID]
		if !ok {
			return fmt.Errorf("diet mapping not found for source diet id=%d (patient=%s)", row.DietID, row.ID)
		}

		if err := s.warehouseRepo.InsertPatient(ctx, tx, transformer.Patient(nextID, row, roomID, dietID)); err != nil {
			return err
		}

		patientMetaBySourceID[row.ID] = patientMeta{
			DWID:         nextID,
			SourceDietID: row.DietID,
			DWRoomID:     roomID,
			DWDietID:     dietID,
		}

		nextID++
	}

	return nil
}

func (s *Service) loadTimes(ctx context.Context, tx *sql.Tx, rows []source.Transaction, timeMap map[string]int64) error {
	seen := make(map[string]struct{})
	unique := make([]source.Transaction, 0, len(rows))

	for _, row := range rows {
		key := utils.TimestampKey(row.OrderDate)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, row)
	}

	sort.SliceStable(unique, func(i, j int) bool {
		if unique[i].OrderDate.Equal(unique[j].OrderDate) {
			return unique[i].ID < unique[j].ID
		}
		return unique[i].OrderDate.Before(unique[j].OrderDate)
	})

	nextID := int64(2)

	for _, row := range unique {
		key := utils.TimestampKey(row.OrderDate)
		timeMap[key] = nextID

		if err := s.warehouseRepo.InsertTime(ctx, tx, transformer.TimeDim(nextID, row.OrderDate)); err != nil {
			return err
		}

		nextID++
	}

	return nil
}

func (s *Service) loadFacts(
	ctx context.Context,
	tx *sql.Tx,
	rows []source.Transaction,
	patientMetaBySourceID map[string]patientMeta,
	menuMapBySourceDiet map[int64]int64,
	timeMap map[string]int64,
) error {
	nextID := int64(1)

	for _, row := range rows {
		meta, ok := patientMetaBySourceID[row.PatientID]
		if !ok {
			return fmt.Errorf("patient mapping not found for source patient id=%s", row.PatientID)
		}

		menuID, ok := menuMapBySourceDiet[meta.SourceDietID]
		if !ok {
			return fmt.Errorf("menu mapping not found for source diet id=%d (patient_id=%s)", meta.SourceDietID, row.PatientID)
		}

		timeID, ok := timeMap[utils.TimestampKey(row.OrderDate)]
		if !ok {
			return fmt.Errorf("time mapping not found for order date=%s", row.OrderDate.String())
		}

		fact := transformer.FactTransaction(
			nextID,
			meta.DWID,
			menuID,
			meta.DWDietID,
			meta.DWRoomID,
			1,
			timeID,
			row.TotalPrice,
			1,
			row.Status,
		)

		if err := s.warehouseRepo.InsertFact(ctx, tx, fact); err != nil {
			return err
		}

		nextID++
	}

	return nil
}