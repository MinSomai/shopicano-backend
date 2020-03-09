package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	LocationTypeCountry LocationType = 0
	LocationTypeState   LocationType = 1
	LocationTypeCity    LocationType = 2
)

type LocationType int64

type Location struct {
	ID          int64        `json:"id" gorm:"column:id;primary_key"`
	Name        string       `json:"name" gorm:"column:name;not null;index"`
	Type        LocationType `json:"-" gorm:"column:type;not null;index"`
	ParentID    int64        `json:"-" gorm:"column:parent_id;type:bigint;not null;index"`
	IsPublished int64        `json:"is_published" gorm:"column:is_published;not null;index"`
	CreatedAt   time.Time    `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (l *Location) TableName() string {
	return "locations"
}

func (l *Location) Populate(db *gorm.DB) error {
	query := `INSERT INTO %s (id, name, TYPE, parent_id, is_published)
		VALUES(1, 'Afghanistan', 0, 0, 1), (2, 'Albania', 0, 0, 1), (3, 'Algeria', 0, 0, 1), (4, 'Andorra', 0, 0, 1), (5, 'Angola', 0, 0, 1), (6, 'Antigua and Barbuda', 0, 0, 1), (7, 'Argentina', 0, 0, 1), (8, 'Armenia', 0, 0, 1), (9, 'Australia', 0, 0, 1), (10, 'Austria', 0, 0, 1), (11, 'Azerbaijan', 0, 0, 1), (12, 'Bahamas', 0, 0, 1), (13, 'Bahrain', 0, 0, 1), (14, 'Bangladesh', 0, 0, 1), (15, 'Barbados', 0, 0, 1), (16, 'Belarus', 0, 0, 1), (17, 'Belgium', 0, 0, 1), (18, 'Belize', 0, 0, 1), (19, 'Benin', 0, 0, 1), (20, 'Bhutan', 0, 0, 1), (21, 'Bolivia', 0, 0, 1), (22, 'Bosnia and Herzegovina', 0, 0, 1), (23, 'Botswana', 0, 0, 1), (24, 'Brazil', 0, 0, 1), (25, 'Brunei ', 0, 0, 1), (26, 'Bulgaria', 0, 0, 1), (27, 'Burkina Faso', 0, 0, 1), (28, 'Burundi', 0, 0, 1), (29, 'Côte dIvoire', 0, 0, 1), (30, 'Cabo Verde', 0, 0, 1), (31, 'Cambodia', 0, 0, 1), (32, 'Cameroon', 0, 0, 1), (33, 'Canada', 0, 0, 1), (34, 'Central African Republic', 0, 0, 1), (35, 'Chad', 0, 0, 1), (36, 'Chile', 0, 0, 1), (37, 'China', 0, 0, 1), (38, 'Colombia', 0, 0, 1), (39, 'Comoros', 0, 0, 1), (40, 'Congo (Congo-Brazzaville)', 0, 0, 1), (41, 'Costa Rica', 0, 0, 1), (42, 'Croatia', 0, 0, 1), (43, 'Cuba', 0, 0, 1), (44, 'Cyprus', 0, 0, 1), (45, 'Czechia (Czech Republic)', 0, 0, 1), (46, 'Democratic Republic of the Congo', 0, 0, 1), (47, 'Denmark', 0, 0, 1), (48, 'Djibouti', 0, 0, 1), (49, 'Dominica', 0, 0, 1), (50, 'Dominican Republic', 0, 0, 1), (51, 'Ecuador', 0, 0, 1), (52, 'Egypt', 0, 0, 1), (53, 'El Salvador', 0, 0, 1), (54, 'Equatorial Guinea', 0, 0, 1), (55, 'Eritrea', 0, 0, 1), (56, 'Estonia', 0, 0, 1), (57, 'Eswatini (fmr. "Swaziland")', 0, 0, 1), (58, 'Ethiopia', 0, 0, 1), (59, 'Fiji', 0, 0, 1), (60, 'Finland', 0, 0, 1), (61, 'France', 0, 0, 1), (62, 'Gabon', 0, 0, 1), (63, 'Gambia', 0, 0, 1), (64, 'Georgia', 0, 0, 1), (65, 'Germany', 0, 0, 1), (66, 'Ghana', 0, 0, 1), (67, 'Greece', 0, 0, 1), (68, 'Grenada', 0, 0, 1), (69, 'Guatemala', 0, 0, 1), (70, 'Guinea', 0, 0, 1), (71, 'Guinea-Bissau', 0, 0, 1), (72, 'Guyana', 0, 0, 1), (73, 'Haiti', 0, 0, 1), (74, 'Holy See', 0, 0, 1), (75, 'Honduras', 0, 0, 1), (76, 'Hungary', 0, 0, 1), (77, 'Iceland', 0, 0, 1), (78, 'India', 0, 0, 1), (79, 'Indonesia', 0, 0, 1), (80, 'Iran', 0, 0, 1), (81, 'Iraq', 0, 0, 1), (82, 'Ireland', 0, 0, 1), (83, 'Israel', 0, 0, 1), (84, 'Italy', 0, 0, 1), (85, 'Jamaica', 0, 0, 1), (86, 'Japan', 0, 0, 1), (87, 'Jordan', 0, 0, 1), (88, 'Kazakhstan', 0, 0, 1), (89, 'Kenya', 0, 0, 1), (90, 'Kiribati', 0, 0, 1), (91, 'Kuwait', 0, 0, 1), (92, 'Kyrgyzstan', 0, 0, 1), (93, 'Laos', 0, 0, 1), (94, 'Latvia', 0, 0, 1), (95, 'Lebanon', 0, 0, 1), (96, 'Lesotho', 0, 0, 1), (97, 'Liberia', 0, 0, 1), (98, 'Libya', 0, 0, 1), (99, 'Liechtenstein', 0, 0, 1), (100, 'Lithuania', 0, 0, 1), (101, 'Luxembourg', 0, 0, 1), (102, 'Madagascar', 0, 0, 1), (103, 'Malawi', 0, 0, 1), (104, 'Malaysia', 0, 0, 1), (105, 'Maldives', 0, 0, 1), (106, 'Mali', 0, 0, 1), (107, 'Malta', 0, 0, 1), (108, 'Marshall Islands', 0, 0, 1), (109, 'Mauritania', 0, 0, 1), (110, 'Mauritius', 0, 0, 1), (111, 'Mexico', 0, 0, 1), (112, 'Micronesia', 0, 0, 1), (113, 'Moldova', 0, 0, 1), (114, 'Monaco', 0, 0, 1), (115, 'Mongolia', 0, 0, 1), (116, 'Montenegro', 0, 0, 1), (117, 'Morocco', 0, 0, 1), (118, 'Mozambique', 0, 0, 1), (119, 'Myanmar (formerly Burma)', 0, 0, 1), (120, 'Namibia', 0, 0, 1), (121, 'Nauru', 0, 0, 1), (122, 'Nepal', 0, 0, 1), (123, 'Netherlands', 0, 0, 1), (124, 'New Zealand', 0, 0, 1), (125, 'Nicaragua', 0, 0, 1), (126, 'Niger', 0, 0, 1), (127, 'Nigeria', 0, 0, 1), (128, 'North Korea', 0, 0, 1), (129, 'North Macedonia', 0, 0, 1), (130, 'Norway', 0, 0, 1), (131, 'Oman', 0, 0, 1), (132, 'Pakistan', 0, 0, 1), (133, 'Palau', 0, 0, 1), (134, 'Palestine State', 0, 0, 1), (135, 'Panama', 0, 0, 1), (136, 'Papua New Guinea', 0, 0, 1), (137, 'Paraguay', 0, 0, 1), (138, 'Peru', 0, 0, 1), (139, 'Philippines', 0, 0, 1), (140, 'Poland', 0, 0, 1), (141, 'Portugal', 0, 0, 1), (142, 'Qatar', 0, 0, 1), (143, 'Romania', 0, 0, 1), (144, 'Russia', 0, 0, 1), (145, 'Rwanda', 0, 0, 1), (146, 'Saint Kitts and Nevis', 0, 0, 1), (147, 'Saint Lucia', 0, 0, 1), (148, 'Saint Vincent and the Grenadines', 0, 0, 1), (149, 'Samoa', 0, 0, 1), (150, 'San Marino', 0, 0, 1), (151, 'Sao Tome and Principe', 0, 0, 1), (152, 'Saudi Arabia', 0, 0, 1), (153, 'Senegal', 0, 0, 1), (154, 'Serbia', 0, 0, 1), (155, 'Seychelles', 0, 0, 1), (156, 'Sierra Leone', 0, 0, 1), (157, 'Singapore', 0, 0, 1), (158, 'Slovakia', 0, 0, 1), (159, 'Slovenia', 0, 0, 1), (160, 'Solomon Islands', 0, 0, 1), (161, 'Somalia', 0, 0, 1), (162, 'South Africa', 0, 0, 1), (163, 'South Korea', 0, 0, 1), (164, 'South Sudan', 0, 0, 1), (165, 'Spain', 0, 0, 1), (166, 'Sri Lanka', 0, 0, 1), (167, 'Sudan', 0, 0, 1), (168, 'Suriname', 0, 0, 1), (169, 'Sweden', 0, 0, 1), (170, 'Switzerland', 0, 0, 1), (171, 'Syria', 0, 0, 1), (172, 'Tajikistan', 0, 0, 1), (173, 'Tanzania', 0, 0, 1), (174, 'Thailand', 0, 0, 1), (175, 'Timor-Leste', 0, 0, 1), (176, 'Togo', 0, 0, 1), (177, 'Tonga', 0, 0, 1), (178, 'Trinidad and Tobago', 0, 0, 1), (179, 'Tunisia', 0, 0, 1), (180, 'Turkey', 0, 0, 1), (181, 'Turkmenistan', 0, 0, 1), (182, 'Tuvalu', 0, 0, 1), (183, 'Uganda', 0, 0, 1), (184, 'Ukraine', 0, 0, 1), (185, 'United Arab Emirates', 0, 0, 1), (186, 'United Kingdom', 0, 0, 1), (187, 'United States of America', 0, 0, 1), (188, 'Uruguay', 0, 0, 1), (189, 'Uzbekistan', 0, 0, 1), (190, 'Vanuatu', 0, 0, 1), (191, 'Venezuela', 0, 0, 1), (192, 'Vietnam', 0, 0, 1), (193, 'Yemen', 0, 0, 1), (194, 'Zambia', 0, 0, 1), (195, 'Zimbabwe', 0, 0, 1)`
	query = fmt.Sprintf(query, l.TableName())
	if err := db.Exec(query).Error; err != nil {
		return err
	}
	return nil
}
