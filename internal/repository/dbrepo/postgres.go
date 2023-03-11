package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Shobhitdimri01/Bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
			end_date, room_id, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restriction (start_date, end_date, room_id, reservation_id,	
			created_at, updated_at, restriction_id) 
			values
			($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	if err != nil {
		return err
	}
	return nil
}

// Will check the date and See whether there is availability in room for roomId --> returns True or false --> if no Availability
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
				select
					count(id)
				from
					room_restriction
				where 
					room_id = $1
					and
					$2 < end_date and $3 > start_date;`

	var numRows int
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	//Check Condition
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityforAll rooms returns a slice of available rooms if any , for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `
				select
					r.id, r.room_name
				from
					rooms r
				where 
					r.id not in 
					(select room_id from room_restriction rr where $1 < rr.end_Date and $2 > rr.start_date ) ;`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room

		err := rows.Scan(
			&room.ID,
			&room.RoomName)

		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}

// GetRoomByID will get room name with its defined Id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
				select id,room_name,created_at,updated_at from rooms where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil

}

// Get User function returns a user by id.
func (m *postgresDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id ,first_name,last_name,email,password,access_level,created_at,updated_at
	from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil

}

// Update User updates the user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name=$1, last_name=$2, email=$3,access_level = $4, updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// Compare Email
func (m *postgresDBRepo) EmailCheck(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Appending an array from database
	query := `SELECT  email FROM public.users;`
	results := make([]string, 0)
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	var scanString string
	for rows.Next() {
		rows.Scan(&scanString)
		results = append(results, scanString)
	}
	m.App.InfoLog.Println("results : ", results)
	for i := 0; i < len(results); i++ {
		if email == results[i] {
			return true
		}
	}
	return false

}

//Hash in the database should correspond to the password by user.
//Authenticate will authemticate the user by matching user data with database

func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	var accesslevel int

	row := m.DB.QueryRowContext(ctx, "select id,password,access_level from users where email = $1", email)
	err := row.Scan(
		&id,
		&hashedPassword,
		&accesslevel,
	)
	if err != nil {
		return id, "", 0, err
	}

	//Comparing user password with Database

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, "", 0, err
	}

	return id, hashedPassword, accesslevel, nil
}

// returns all reservation from database
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
		r.end_date, r.room_id, r.created_at, r.updated_at,r.processed,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		order by r.start_date asc
`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}

// returns AllNewreservation from database
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
		r.end_date, r.room_id, r.created_at, r.updated_at,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where processed = 0
		order by r.start_date asc
`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}
func (m *postgresDBRepo) GetAdminByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var admin models.User
	query := `SELECT id, first_name, last_name, email, access_level FROM public.users WHERE id=$1;`
	rows := m.DB.QueryRowContext(ctx, query, id)
	err := rows.Scan(
		&admin.ID,
		&admin.FirstName,
		&admin.LastName,
		&admin.Email,
		&admin.AccessLevel,
	)
	if err != nil {
		return admin, err
	}
	fmt.Println("All-Data:", admin)
	return admin, nil
}
func (m *postgresDBRepo) UpdateAdminData(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	fmt.Println("Hellooooooooooooooooooooooo", u.ID)
	query := `UPDATE public.users
	SET first_name=$1, last_name=$2, email=$3, access_level=$4 WHERE id=$5;`
	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteAdminByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `DELETE FROM public.users WHERE id=$1;`
	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}
func (m *postgresDBRepo) GetAllAdmins() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var admin []models.User
	query := `SELECT u.id, u.first_name, u.last_name, u.email, u.access_level  
	FROM public.users u;`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return admin, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.User
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.AccessLevel,
		)
		if err != nil {
			return admin, err
		}
		fmt.Println(i)
		admin = append(admin, i)
	}
	if err = rows.Err(); err != nil {
		return admin, err
	}

	return admin, nil

}

func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
		r.end_date, r.room_id, r.created_at, r.updated_at,r.processed,
		rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where r.id = $1
`
	rows := m.DB.QueryRowContext(ctx, query, id)
	err := rows.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set first_name=$1, last_name=$2, email=$3, phone = $4, updated_at = $5
	where id = $6
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation delete reservations by id
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id=$1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

// updates process for reservation by id
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set processed=$1 where id=$2`

	_, err := m.DB.ExecContext(ctx, query, processed, id)

	if err != nil {
		return err
	}

	return nil
}

//Counting total IDs from database

func (m *postgresDBRepo) CountReservation() (total_res, del_res, sun_res int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select count(id) from reservations;`
	query1 := `select count(id) from reservations where room_id = 1;`
	query2 := `select count(id) from reservations where room_id = 2;`
	row := m.DB.QueryRowContext(ctx, query)
	row1 := m.DB.QueryRowContext(ctx, query1)
	row2 := m.DB.QueryRowContext(ctx, query2)
	err := row.Scan(&total_res)
	res2 := row1.Scan(&del_res)
	res3 := row2.Scan(&sun_res)
	if err != nil || res2 != nil || res3 != nil {
		fmt.Print(err)
	}
	return total_res, del_res, sun_res
}

func (m *postgresDBRepo) CountMonths()([]int,[]int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var Month []int
	var BookingCount []int
	query := `SELECT EXTRACT('MONTH' FROM start_date) AS month,
	COUNT(id) AS data
	FROM reservations
	GROUP BY EXTRACT('MONTH' FROM start_date) order by month asc;`
	rows, _ := m.DB.QueryContext(ctx, query)
	for rows.Next(){
		var a int
		var b int

		err := rows.Scan(
			&a,
			&b,
		)
		if err!= nil{
			fmt.Println(err)
		}
		Month = append(Month, a)
		BookingCount = append(BookingCount, b)
	}
	m.App.InfoLog.Println("Month:",Month,"Booking:",BookingCount)
	return Month,BookingCount
	
}
func (m *postgresDBRepo) InsertUserData(r models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO public.users
	(first_name, last_name, email, "password", access_level, created_at, updated_at)
	VALUES($1,$2,$3,$4,$5,$6,$7);`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Password,
		r.AccessLevel,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err

	}
	return nil
}
func (m *postgresDBRepo)AllRooms()([]models.Room,error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `Select id,room_name,created_at,updated_at from rooms order by room_name`
	row,err :=m.DB.QueryContext(ctx,query)
	if err!=nil{
		return rooms,err
	}
	defer row.Close()

	for row.Next(){
		var i models.Room
		err := row.Scan(
			&i.ID,
			&i.RoomName,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err!= nil{
			return rooms,err
		}
		rooms = append(rooms, i)
	}
	if err!=nil{
		return rooms,err
	}
	return rooms,nil
}


// GetRestrictionsForRoomByDate returns restrictions for a room by date range
func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, startDate, endDate time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction
	fmt.Println("Helloooo")
	query := `
		select id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date
		from room_restriction where $1 < end_date and $2 >= start_date
		and room_id = $3
`

	rows, err := m.DB.QueryContext(ctx, query,startDate, endDate,roomID)
	if err != nil {
		return nil, err
	}
	fmt.Println("rows",rows)
	defer rows.Close()
	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		fmt.Println("Myerr",err)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, r)
	}
	fmt.Println("restriction",restrictions)
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restrictions, nil
}
// InsertBlockForRoom inserts a room restriction
func (m *postgresDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restriction (start_date, end_date, room_id, restriction_id,
			created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// DeleteBlockByID deletes a room restriction
func (m *postgresDBRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from room_restriction where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}