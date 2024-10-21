package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iki-rumondor/go-speech/internal/config"
	"github.com/iki-rumondor/go-speech/internal/middleware"
)

func StartServer(handlers *config.Handlers) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:5173", "http://103.26.13.166:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12,
	}))

	public := router.Group("api")
	{
		public.POST("/user", handlers.UserHandler.CreateUser)
		public.GET("/user/:uuid", handlers.UserHandler.GetUser)
		public.POST("/verify-user", handlers.UserHandler.VerifyUser)
		public.GET("/public/classes", handlers.UserHandler.GetAllClasses)
		public.GET("/roles", handlers.UserHandler.GetRoles)
		public.GET("/department", handlers.MasterHandler.GetAllDepartment)
		public.GET("/file/videos/:videoName", handlers.FileHandler.GetFileVideo)
		public.GET("/file/subtitle/:subtitle", handlers.FileHandler.GetFileSubtitle)
		public.GET("/file/books/:bookName", handlers.FileHandler.GetFileBook)
		public.GET("/file/answers/:filename", handlers.FileHandler.GetFileAnswer)
	}

	user := router.Group("api").Use(middleware.IsValidJWT())
	{
		user.GET("/videos/:uuid", handlers.FileHandler.GetVideo)
		user.GET("/videos/classes/:uuid", handlers.FileHandler.GetClassVideos)
		user.GET("/classes/:uuid", handlers.MasterHandler.GetClass)
		user.GET("/books/classes/:uuid", handlers.FileHandler.GetClassBooks)
		user.GET("/books/:uuid", handlers.FileHandler.GetVideo)
		user.GET("/notes/classes/:uuid", handlers.MasterHandler.GetNotes)
		user.GET("/materials/classes/:class_uuid", handlers.MasterHandler.GetAllMaterials)
		user.GET("/materials/:uuid", handlers.MasterHandler.GetMaterial)
		user.GET("/classes/all", handlers.MasterHandler.GetAllClasses)
		user.GET("/subjects/all", handlers.MasterHandler.GetAllSubjects)
	}

	admin := router.Group("api").Use(middleware.IsValidJWT(), middleware.IsRole("ADMIN"))
	{
		admin.POST("/department", handlers.MasterHandler.CreateDepartment)
		admin.GET("/department/:uuid", handlers.MasterHandler.GetDepartment)
		admin.PUT("/department/:uuid", handlers.MasterHandler.UpdateDepartment)
		admin.DELETE("/department/:uuid", handlers.MasterHandler.DeleteDepartment)

		admin.POST("/teachers/import", handlers.ImportHandler.ImportTeachers)
		admin.POST("/teachers", handlers.MasterHandler.CreateTeacher)
		admin.GET("/teachers", handlers.UserHandler.GetTeachers)
		admin.GET("/teachers/:uuid", handlers.MasterHandler.GetTeacher)
		admin.PUT("/teachers/:uuid", handlers.MasterHandler.UpdateTeacher)
		admin.DELETE("/teachers/:uuid", handlers.MasterHandler.DeleteTeacher)
		admin.PATCH("/teacher/:uuid/activate", handlers.UserHandler.ActivateUser)

		admin.GET("/students", handlers.MasterHandler.GetStudents)
		admin.GET("/students/:uuid", handlers.MasterHandler.GetStudent)
		admin.POST("/students/import", handlers.ImportHandler.ImportStudents)
		admin.PUT("/students/:uuid", handlers.MasterHandler.UpdateStudent)
		admin.DELETE("/students/:uuid", handlers.MasterHandler.DeleteStudent)

		admin.GET("dashboards/admin", handlers.UserHandler.DashboardAdmin)

		admin.POST("/pdf/reports/classes/:uuid", handlers.MasterHandler.GetClassesReport)
		admin.GET("/assignments/students/:studentUuid", handlers.AssignmentHandler.FindAssignmentByStudent)
		admin.POST("/classes", handlers.MasterHandler.CreateClass)
		admin.POST("/subjects", handlers.MasterHandler.CreateSubject)
		admin.PUT("/classes/:uuid", handlers.MasterHandler.UpdateClass)
		admin.DELETE("/classes/:uuid", handlers.MasterHandler.DeleteClass)
		admin.GET("/classes/request", handlers.UserHandler.GetRequestClasses)
		admin.PATCH("/classes/:uuid/request", handlers.UserHandler.UpdateStatusClassReq)
	}

	teacher := router.Group("api").Use(middleware.IsValidJWT(), middleware.IsRole("DOSEN"), middleware.SetUserUuid())
	{
		teacher.GET("/classes", handlers.MasterHandler.GetTeacherClasses)
		teacher.POST("/videos", handlers.FileHandler.CreateVideo)
		teacher.PUT("/videos/:uuid", handlers.FileHandler.UpdateVideo)
		teacher.DELETE("/videos/:uuid", handlers.FileHandler.DeleteVideo)

		teacher.POST("/books", handlers.FileHandler.CreateBook)
		teacher.DELETE("/books/:uuid", handlers.FileHandler.DeleteBook)

		teacher.POST("/notes", handlers.MasterHandler.CreateNote)
		teacher.GET("/notes/:uuid", handlers.MasterHandler.GetNote)
		teacher.PUT("/notes/:uuid", handlers.MasterHandler.UpdateNote)
		teacher.DELETE("/notes/:uuid", handlers.MasterHandler.DeleteNote)

		teacher.GET("dashboards/teacher", handlers.UserHandler.DashboardTeacher)
		teacher.POST("/materials", handlers.MasterHandler.CreateMaterial)
		teacher.PUT("/materials/:uuid", handlers.MasterHandler.UpdateMaterial)
		teacher.DELETE("/materials/:uuid", handlers.MasterHandler.DeleteMaterial)

		teacher.POST("/assignments", handlers.AssignmentHandler.CreateAssignment)
		teacher.GET("/assignments/classes/:classUuid", handlers.AssignmentHandler.FindAssignmentByClass)
		teacher.GET("/assignments/:uuid", handlers.AssignmentHandler.FirstAssignmentByUuid)
		teacher.PUT("/assignments/:uuid", handlers.AssignmentHandler.UpdateAssignment)
		teacher.DELETE("/assignments/:uuid", handlers.AssignmentHandler.DeleteAssignment)

		teacher.PATCH("/answers/:answerUuid/grading", handlers.AssignmentHandler.GradeAnswer)
		teacher.GET("/students/classes/:classUuid", handlers.UserHandler.GetStudentsByClass)
	}

	student := router.Group("api").Use(middleware.IsValidJWT(), middleware.IsRole("MAHASISWA"), middleware.SetUserUuid())
	{
		student.GET("/informations/students", handlers.UserHandler.GetStudentInformation)

		student.POST("/class/register", handlers.UserHandler.CreateClassRequest)
		student.POST("/classes/join", handlers.UserHandler.JoinClass)
		student.GET("/class/request/students", handlers.UserHandler.GetStudentRequestClasses)
		student.GET("/classes/students/:userUuid", handlers.UserHandler.GetStudentClasses)

		student.GET("dashboards/student", handlers.UserHandler.DashboardStudent)
		student.GET("/assignments/students/classes/:classUuid", handlers.AssignmentHandler.FindAssignmentByUser)
		student.POST("/answers/assignments/:assignmentUuid", handlers.AssignmentHandler.UploadAnswer)
		student.GET("/notifications", handlers.UserHandler.GetNotifications)
		student.POST("/notifications/read", handlers.UserHandler.ReadNotification)
	}

	return router
}
