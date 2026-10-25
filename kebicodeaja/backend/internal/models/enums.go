package models

type UserRole string
type CTAspect string
type ChallengeType string
type Difficulty string
type AttachmentType string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"

	CTDecomposition      CTAspect = "decomposition"
	CTPatternRecognition CTAspect = "pattern_recognition"
	CTAbstraction        CTAspect = "abstraction"
	CTAlgorithm          CTAspect = "algorithm"

	ChallengeMultipleChoice ChallengeType = "multiple_choice"
	ChallengeOrdering       ChallengeType = "ordering"
	ChallengeFillBlank      ChallengeType = "fill_blank"

	DiffEasy   Difficulty = "easy"
	DiffMedium Difficulty = "medium"
	DiffHard   Difficulty = "hard"

	AttachVideo AttachmentType = "video"
	AttachEbook AttachmentType = "ebook"
	AttachFile  AttachmentType = "file"
	AttachLink  AttachmentType = "link"
)
