package dbinstance

import "github.com/kloeckner-i/db-operator/pkg/test"

func testAmazonMysqlInstance() *Amazon {
	return &Amazon{Generic: Generic{Host: test.GetMysqlHost(), Port: test.GetMysqlPort(), Engine: "mysql", User: "root", Password: test.GetMysqlAdminPassword()}}
}

func testAmazonPostgresInstance() *Amazon {
	return &Amazon{Generic: Generic{Host: test.GetPostgresHost(), Port: test.GetPostgresPort(), Engine: "postgres", User: "postgres", Password: test.GetPostgresAdminPassword()}}
}
