package main

func main() {
	a := App{}
	//a.Initialize(os.Getenv("APP_DB_NAME"))
	a.Initialize("test_database.db")

	a.Run(":8080")
}
