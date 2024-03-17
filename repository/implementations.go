package repository

import "context"

func (r *Repository) InsertUsersData(ctx context.Context, input InsertUsersDataInput) (output InsertUsersDataOutput, err error) {
	stmt, err := r.Db.PrepareContext(ctx, "INSERT INTO users(phone_number, full_name, hashed_pass, count) VALUES($1, $2, $3, 0) RETURNING id")
	if err != nil {
		return output, err
	}
	defer stmt.Close()

	var id int32
	err = stmt.QueryRowContext(ctx, input.PhoneNumber, input.FullName, input.HashedPassword).Scan(&id)
	if err != nil {
		return output, err
	}

	output.Id = id
	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, input GetUserByPhoneNumberInput) (output GetUserByPhoneNumberOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, hashed_pass, count FROM users WHERE phone_number=$1", input.PhoneNumber).Scan(&output.Id, &output.FullName, &output.PhoneNumber, &output.HashedPass, &output.Count)

	return
}

func (r *Repository) GetUserById(ctx context.Context, input GetUserByIdInput) (output GetUserByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, hashed_pass, count FROM users WHERE id=$1", input.Id).Scan(&output.Id, &output.FullName, &output.PhoneNumber, &output.HashedPass, &output.Count)

	return
}

func (r *Repository) UpdateFullNameByIdInput(ctx context.Context, input UpdateFullNameByIdInput) (err error) {
	stmt, err := r.Db.PrepareContext(ctx, "UPDATE users SET full_name=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, input.FullName, input.Id)

	return
}

func (r *Repository) UpdateCountById(ctx context.Context, input UpdateCountByIdInput) (err error) {
	stmt, err := r.Db.PrepareContext(ctx, "UPDATE users SET count=$1 WHERE id=$2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, input.Count, input.Id)

	return
}
