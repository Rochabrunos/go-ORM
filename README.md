# DVDRENTAL CRUD

This repository contains an application made with the database of the DVDRental example from the PostgreSQLTutorial site.

The main goal of this project is to get knowledge of ORM tools and Web Framework available for GO Lang. In this sense, I had to choose which library/framework to use for developing the project. The decision was based on the popularity of the tool. Finally, the conclusion leads to the following set of libraries:

* [GORM](https://gorm.io/docs/): it is a full-featured ORM library
* [GIN](https://gin-gonic.com/docs/introduction/): an HTTP framework written in GO Lang

Make sure to include the environment variables for running the project: DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT.

## ENDPOINTS:
- LANGUAGE:
	- GET "/languages?p=1" : get all Languages
		- The resource is limited to 10 entries
			- "p" stands for page, e.g. p=0 returns the set of id [1-10], for p=1 [11-20], ... . (default p=0)
		- Retruns a JSON of Languages
	- GET "/languages/:id" : get Language by id
		- The :id must be a valid one, otherwise an error will be returned
		- Returns a JSON Language with the given id
	- POST "/languages" : create a Language
		- The body should contain the "name" attribute for Language in JSON format
		- Returns a JSON with the created Language
	- PUT "/languages/:id" : modify Language by id
		- The :id must be a valid one, otherwise an error will be returned
		- The body should contain the "name" attribute for Language in JSON format
		- Returns a JSON with the changed Language
- CATEGORY:
	- GET "/categories?p=1" : get all Categories
		- The resource is limited to 10 entries
			- "p" stands for page, e.g. p=0 returns the set of id [1-10], for p=1 [11-20], ... . (default p=0)
		- Returns a JSON object with all Categories
	- GET "/categories/:id" : get Category by id
		- The :id must be a valid one, otherwise an error will be returned
		- Returns a JSON with Category with the given id
	- POST "/categories" : create a Category
		- The body should contain the "name" attribute for Category in JSON format
		- Returns a JSON with the created Category
	- PUT "/categories/:id" : modify category by id
		- The :id must be a valid one, otherwise an error will be returned
		- The body should contain the "name" attribute for Category in JSON format
		- Returns a JSON with the changed Category
