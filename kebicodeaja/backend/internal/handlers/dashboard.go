package handlers

import (
	"log"
	"math"
	"net/http"
	"time"

	"github.com/daffafauzan4445/kebcodeaja/backend/internal/config"
	"github.com/gin-gonic/gin"
)

var ctAspects = []string{"decomposition", "pattern_recognition", "abstraction", "algorithm"}

type aspectValue struct {
	Aspect string `json:"aspect"`
	Value  int    `json:"value"`
}

type aspectRow struct {
	Aspect   string  `gorm:"column:aspect"`
	AvgValue float64 `gorm:"column:avg_value"`
}

func fillAspects(rows []aspectRow) []aspectValue {
	byAspect := make(map[string]int, len(rows))
	for _, r := range rows {
		byAspect[r.Aspect] = int(math.Round(r.AvgValue))
	}
	result := make([]aspectValue, 0, len(ctAspects))
	for _, a := range ctAspects {
		result = append(result, aspectValue{Aspect: a, Value: byAspect[a]})
	}
	return result
}

type teacherClassRow struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	ClassCode      string `json:"class_code"`
	StudentCount   int    `json:"student_count"`
	MaterialCount  int    `json:"material_count"`
	ChallengeCount int    `json:"challenge_count"`
}

type teacherSubmissionRow struct {
	ID             uint      `json:"id"`
	StudentName    string    `json:"student_name"`
	ChallengeTitle string    `json:"challenge_title"`
	ClassName      string    `json:"class_name"`
	Score          int       `json:"score"`
	SubmittedAt    time.Time `json:"submitted_at"`
}

func TeacherDashboard(c *gin.Context) {
	teacherID := c.GetUint("user_id")
	db := config.DB

	classes := []teacherClassRow{}
	err := db.Raw(`
		SELECT c.id, c.name, c.class_code,
			(SELECT COUNT(*) FROM class_members cm WHERE cm.class_id = c.id) AS student_count,
			(SELECT COUNT(*) FROM materials m WHERE m.class_id = c.id) AS material_count,
			(SELECT COUNT(*) FROM challenges ch WHERE ch.class_id = c.id) AS challenge_count
		FROM classes c
		WHERE c.teacher_id = ?
		ORDER BY c.created_at DESC`, teacherID).Scan(&classes).Error
	if err != nil {
		log.Printf("dashboard guru, data kelas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data kelas", "error": debugError(err)})
		return
	}
	if classes == nil {
		classes = []teacherClassRow{}
	}

	var totalStudents int64
	err = db.Raw(`
		SELECT COUNT(DISTINCT cm.student_id)
		FROM class_members cm
		JOIN classes c ON c.id = cm.class_id
		WHERE c.teacher_id = ?`, teacherID).Scan(&totalStudents).Error
	if err != nil {
		log.Printf("dashboard guru, data siswa: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data siswa", "error": debugError(err)})
		return
	}

	totalMaterials := 0
	totalChallenges := 0
	for _, cl := range classes {
		totalMaterials += cl.MaterialCount
		totalChallenges += cl.ChallengeCount
	}

	var avgRows []aspectRow
	err = db.Raw(`
		SELECT ch.ct_aspect AS aspect,
			COALESCE(AVG(s.score * 100.0 / NULLIF(ch.max_score, 0)), 0) AS avg_value
		FROM submissions s
		JOIN challenges ch ON ch.id = s.challenge_id
		JOIN classes c ON c.id = ch.class_id
		WHERE c.teacher_id = ?
		GROUP BY ch.ct_aspect`, teacherID).Scan(&avgRows).Error
	if err != nil {
		log.Printf("dashboard guru, rata-rata nilai: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat rata-rata nilai", "error": debugError(err)})
		return
	}

	submissions := []teacherSubmissionRow{}
	err = db.Raw(`
		SELECT s.id, u.name AS student_name, ch.title AS challenge_title,
			c.name AS class_name, s.score, s.submitted_at
		FROM submissions s
		JOIN challenges ch ON ch.id = s.challenge_id
		JOIN classes c ON c.id = ch.class_id
		JOIN users u ON u.id = s.student_id
		WHERE c.teacher_id = ?
		ORDER BY s.submitted_at DESC
		LIMIT 5`, teacherID).Scan(&submissions).Error
	if err != nil {
		log.Printf("dashboard guru, pengumpulan terbaru: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat pengumpulan terbaru", "error": debugError(err)})
		return
	}
	if submissions == nil {
		submissions = []teacherSubmissionRow{}
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"classes":    len(classes),
			"students":   totalStudents,
			"materials":  totalMaterials,
			"challenges": totalChallenges,
		},
		"classes":            classes,
		"averages":           fillAspects(avgRows),
		"recent_submissions": submissions,
	})
}

type studentChallengeRow struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	CTAspect   string `json:"ct_aspect" gorm:"column:ct_aspect"`
	Difficulty string `json:"difficulty"`
	ClassName  string `json:"class_name"`
	Score      *int   `json:"score"`
}

type studentMaterialRow struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	ClassName string    `json:"class_name"`
	CreatedAt time.Time `json:"created_at"`
}

func StudentDashboard(c *gin.Context) {
	studentID := c.GetUint("user_id")
	db := config.DB

	var classCount int64
	err := db.Raw(`SELECT COUNT(*) FROM class_members WHERE student_id = ?`, studentID).Scan(&classCount).Error
	if err != nil {
		log.Printf("dashboard siswa, data kelas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data kelas", "error": debugError(err)})
		return
	}

	var materialCount int64
	err = db.Raw(`
		SELECT COUNT(*)
		FROM materials m
		JOIN class_members cm ON cm.class_id = m.class_id
		WHERE cm.student_id = ?`, studentID).Scan(&materialCount).Error
	if err != nil {
		log.Printf("dashboard siswa, data materi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data materi", "error": debugError(err)})
		return
	}

	var completed int64
	err = db.Raw(`
		SELECT COUNT(DISTINCT challenge_id)
		FROM submissions
		WHERE student_id = ?`, studentID).Scan(&completed).Error
	if err != nil {
		log.Printf("dashboard siswa, data tantangan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data tantangan", "error": debugError(err)})
		return
	}

	var avgScore float64
	err = db.Raw(`
		SELECT COALESCE(AVG(t.best), 0)
		FROM (
			SELECT MAX(s.score * 100.0 / NULLIF(ch.max_score, 0)) AS best
			FROM submissions s
			JOIN challenges ch ON ch.id = s.challenge_id
			WHERE s.student_id = ?
			GROUP BY s.challenge_id
		) t`, studentID).Scan(&avgScore).Error
	if err != nil {
		log.Printf("dashboard siswa, rata-rata nilai: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat rata-rata nilai", "error": debugError(err)})
		return
	}

	var progressRows []aspectRow
	err = db.Raw(`
		SELECT t.aspect AS aspect, COALESCE(AVG(t.best), 0) AS avg_value
		FROM (
			SELECT ch.ct_aspect AS aspect,
				MAX(s.score * 100.0 / NULLIF(ch.max_score, 0)) AS best
			FROM submissions s
			JOIN challenges ch ON ch.id = s.challenge_id
			WHERE s.student_id = ?
			GROUP BY s.challenge_id, ch.ct_aspect
		) t
		GROUP BY t.aspect`, studentID).Scan(&progressRows).Error
	if err != nil {
		log.Printf("dashboard siswa, progres: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat progres", "error": debugError(err)})
		return
	}

	challenges := []studentChallengeRow{}
	err = db.Raw(`
		SELECT ch.id, ch.title, ch.ct_aspect, ch.difficulty, c.name AS class_name,
			(SELECT MAX(s.score) FROM submissions s
				WHERE s.challenge_id = ch.id AND s.student_id = ?) AS score
		FROM challenges ch
		JOIN classes c ON c.id = ch.class_id
		JOIN class_members cm ON cm.class_id = c.id
		WHERE cm.student_id = ?
		ORDER BY ch.created_at DESC
		LIMIT 5`, studentID, studentID).Scan(&challenges).Error
	if err != nil {
		log.Printf("dashboard siswa, daftar tantangan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat daftar tantangan", "error": debugError(err)})
		return
	}
	if challenges == nil {
		challenges = []studentChallengeRow{}
	}

	materials := []studentMaterialRow{}
	err = db.Raw(`
		SELECT m.id, m.title, c.name AS class_name, m.created_at
		FROM materials m
		JOIN classes c ON c.id = m.class_id
		JOIN class_members cm ON cm.class_id = m.class_id
		WHERE cm.student_id = ?
		ORDER BY m.created_at DESC
		LIMIT 3`, studentID).Scan(&materials).Error
	if err != nil {
		log.Printf("dashboard siswa, daftar materi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat daftar materi", "error": debugError(err)})
		return
	}
	if materials == nil {
		materials = []studentMaterialRow{}
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"classes":              classCount,
			"materials":            materialCount,
			"completed_challenges": completed,
			"average_score":        int(math.Round(avgScore)),
		},
		"progress":   fillAspects(progressRows),
		"challenges": challenges,
		"materials":  materials,
	})
}
