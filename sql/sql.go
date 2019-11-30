package sql

func Parse(sql string) (*Command, error) {

	parsed, err := textParse(sql)

	if err != nil {
		return nil, err
	}

	command, err := doCommand(parsed)

	if err != nil {
		return nil, err
	}

	return command, nil
}
