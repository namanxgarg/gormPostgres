**What I learn?**

GORM --> ORM for Go
ORM?
  interact with go using object, not having to use sql for queries for CRUD
  
log?
  logs it either in terminal or in output.log, also terminates the program if called upon

& --> when not want to create a copy, actual value access
  cuz if you send a copy, GORM won't be able to make changes in actual record

auto migrate --> checks for struct and database schema correspondance in database

*gorm.db --> database connection instance, used in many parts --> more to explore it

if struct function (func (r *structure_name))--> method

fiber --> really fast, based on fasthttp, express.js inspired
  net/http--> base of everything
  gin--> slower than fiber in some concurrent tasks
  gorilla mux --> based on net/http, slower than fiber

context --> request and responses are in context

why postgre --> ACID compliant, supports json --> need to learn a lot more (read a bit, but learn about database)
