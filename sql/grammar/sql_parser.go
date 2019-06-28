// Code generated from /Users/dupeng/Documents/gohub/src/github.com/eleme/lindb/sql/antlr4/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // SQL

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 94, 694,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9,
	39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44, 9, 44,
	4, 45, 9, 45, 4, 46, 9, 46, 4, 47, 9, 47, 4, 48, 9, 48, 4, 49, 9, 49, 4,
	50, 9, 50, 4, 51, 9, 51, 4, 52, 9, 52, 4, 53, 9, 53, 4, 54, 9, 54, 4, 55,
	9, 55, 4, 56, 9, 56, 4, 57, 9, 57, 4, 58, 9, 58, 4, 59, 9, 59, 4, 60, 9,
	60, 4, 61, 9, 61, 4, 62, 9, 62, 4, 63, 9, 63, 4, 64, 9, 64, 4, 65, 9, 65,
	4, 66, 9, 66, 4, 67, 9, 67, 4, 68, 9, 68, 4, 69, 9, 69, 4, 70, 9, 70, 4,
	71, 9, 71, 4, 72, 9, 72, 4, 73, 9, 73, 4, 74, 9, 74, 4, 75, 9, 75, 4, 76,
	9, 76, 4, 77, 9, 77, 4, 78, 9, 78, 4, 79, 9, 79, 4, 80, 9, 80, 4, 81, 9,
	81, 4, 82, 9, 82, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 183, 10, 3,
	3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 5, 4, 190, 10, 4, 3, 4, 3, 4, 5, 4, 194,
	10, 4, 3, 5, 3, 5, 3, 5, 7, 5, 199, 10, 5, 12, 5, 14, 5, 202, 11, 5, 3,
	6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3,
	6, 3, 6, 5, 6, 218, 10, 6, 3, 7, 3, 7, 3, 7, 7, 7, 223, 10, 7, 12, 7, 14,
	7, 226, 11, 7, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3,
	8, 3, 8, 3, 9, 3, 9, 3, 10, 3, 10, 3, 11, 3, 11, 3, 12, 3, 12, 3, 13, 3,
	13, 3, 14, 3, 14, 3, 15, 3, 15, 3, 16, 3, 16, 3, 17, 3, 17, 3, 17, 3, 17,
	3, 17, 5, 17, 260, 10, 17, 3, 17, 3, 17, 5, 17, 264, 10, 17, 3, 18, 3,
	18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 3, 20, 3, 20, 3, 20, 3, 21, 3, 21,
	3, 21, 5, 21, 279, 10, 21, 3, 21, 5, 21, 282, 10, 21, 3, 22, 3, 22, 3,
	22, 3, 22, 3, 22, 3, 22, 5, 22, 290, 10, 22, 3, 23, 3, 23, 3, 23, 3, 23,
	3, 23, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 5, 24, 304, 10,
	24, 3, 24, 5, 24, 307, 10, 24, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25,
	3, 25, 3, 25, 3, 25, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 5, 26, 324,
	10, 26, 3, 27, 3, 27, 3, 27, 5, 27, 329, 10, 27, 3, 28, 3, 28, 3, 28, 3,
	28, 5, 28, 335, 10, 28, 3, 28, 3, 28, 5, 28, 339, 10, 28, 3, 29, 3, 29,
	3, 29, 3, 29, 3, 29, 3, 29, 5, 29, 347, 10, 29, 3, 30, 3, 30, 3, 30, 3,
	30, 3, 30, 3, 31, 3, 31, 3, 31, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 5, 32,
	362, 10, 32, 3, 33, 3, 33, 3, 34, 3, 34, 3, 35, 3, 35, 3, 36, 3, 36, 3,
	37, 5, 37, 373, 10, 37, 3, 37, 3, 37, 3, 37, 3, 37, 5, 37, 379, 10, 37,
	3, 37, 5, 37, 382, 10, 37, 3, 37, 5, 37, 385, 10, 37, 3, 37, 5, 37, 388,
	10, 37, 3, 37, 5, 37, 391, 10, 37, 3, 37, 5, 37, 394, 10, 37, 3, 38, 3,
	38, 3, 38, 7, 38, 399, 10, 38, 12, 38, 14, 38, 402, 11, 38, 3, 39, 3, 39,
	5, 39, 406, 10, 39, 3, 40, 3, 40, 3, 40, 3, 41, 3, 41, 3, 41, 3, 42, 3,
	42, 3, 42, 3, 43, 3, 43, 3, 43, 5, 43, 420, 10, 43, 3, 43, 3, 43, 3, 43,
	7, 43, 425, 10, 43, 12, 43, 14, 43, 428, 11, 43, 3, 44, 3, 44, 3, 44, 3,
	44, 3, 44, 5, 44, 435, 10, 44, 5, 44, 437, 10, 44, 3, 45, 3, 45, 3, 45,
	3, 45, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3,
	46, 3, 46, 3, 46, 3, 46, 5, 46, 456, 10, 46, 3, 46, 3, 46, 3, 46, 3, 46,
	5, 46, 462, 10, 46, 3, 46, 3, 46, 3, 46, 7, 46, 467, 10, 46, 12, 46, 14,
	46, 470, 11, 46, 3, 47, 3, 47, 3, 47, 7, 47, 475, 10, 47, 12, 47, 14, 47,
	478, 11, 47, 3, 48, 3, 48, 3, 48, 5, 48, 483, 10, 48, 3, 49, 3, 49, 3,
	49, 3, 49, 5, 49, 489, 10, 49, 3, 50, 3, 50, 5, 50, 493, 10, 50, 3, 51,
	3, 51, 3, 51, 5, 51, 498, 10, 51, 3, 51, 3, 51, 3, 52, 3, 52, 3, 52, 3,
	52, 3, 52, 3, 52, 3, 52, 3, 52, 5, 52, 510, 10, 52, 3, 52, 5, 52, 513,
	10, 52, 3, 53, 3, 53, 3, 53, 7, 53, 518, 10, 53, 12, 53, 14, 53, 521, 11,
	53, 3, 54, 3, 54, 3, 54, 3, 54, 3, 54, 3, 54, 5, 54, 529, 10, 54, 3, 55,
	3, 55, 3, 56, 3, 56, 3, 56, 3, 56, 3, 57, 3, 57, 3, 57, 3, 57, 3, 58, 3,
	58, 7, 58, 543, 10, 58, 12, 58, 14, 58, 546, 11, 58, 3, 59, 3, 59, 3, 59,
	7, 59, 551, 10, 59, 12, 59, 14, 59, 554, 11, 59, 3, 60, 3, 60, 3, 60, 3,
	61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 5, 61, 565, 10, 61, 3, 61, 3, 61,
	3, 61, 3, 61, 7, 61, 571, 10, 61, 12, 61, 14, 61, 574, 11, 61, 3, 62, 3,
	62, 3, 63, 3, 63, 3, 64, 3, 64, 3, 64, 3, 64, 3, 65, 3, 65, 3, 65, 3, 65,
	3, 65, 3, 65, 3, 65, 3, 65, 5, 65, 592, 10, 65, 3, 66, 3, 66, 3, 66, 3,
	66, 3, 66, 3, 66, 3, 66, 3, 66, 5, 66, 602, 10, 66, 3, 66, 3, 66, 3, 66,
	3, 66, 3, 66, 3, 66, 3, 66, 3, 66, 3, 66, 3, 66, 3, 66, 3, 66, 7, 66, 616,
	10, 66, 12, 66, 14, 66, 619, 11, 66, 3, 67, 3, 67, 3, 67, 3, 68, 3, 68,
	3, 69, 3, 69, 3, 69, 5, 69, 629, 10, 69, 3, 69, 3, 69, 3, 70, 3, 70, 3,
	70, 7, 70, 636, 10, 70, 12, 70, 14, 70, 639, 11, 70, 3, 71, 3, 71, 5, 71,
	643, 10, 71, 3, 72, 3, 72, 5, 72, 647, 10, 72, 3, 72, 3, 72, 5, 72, 651,
	10, 72, 3, 73, 3, 73, 3, 73, 3, 73, 3, 74, 5, 74, 658, 10, 74, 3, 74, 3,
	74, 3, 75, 5, 75, 663, 10, 75, 3, 75, 3, 75, 3, 76, 3, 76, 3, 76, 3, 77,
	3, 77, 3, 78, 3, 78, 3, 79, 3, 79, 3, 80, 3, 80, 3, 81, 3, 81, 5, 81, 680,
	10, 81, 3, 81, 3, 81, 3, 81, 5, 81, 685, 10, 81, 7, 81, 687, 10, 81, 12,
	81, 14, 81, 690, 11, 81, 3, 82, 3, 82, 3, 82, 2, 6, 84, 90, 120, 130, 83,
	2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38,
	40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74,
	76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 108,
	110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 130, 132, 134, 136, 138,
	140, 142, 144, 146, 148, 150, 152, 154, 156, 158, 160, 162, 2, 10, 5, 2,
	48, 48, 71, 73, 78, 78, 3, 2, 40, 41, 4, 2, 43, 44, 92, 93, 3, 2, 46, 47,
	4, 2, 48, 48, 78, 78, 3, 2, 62, 68, 3, 2, 86, 87, 11, 2, 3, 3, 7, 7, 9,
	11, 15, 24, 26, 29, 31, 35, 38, 52, 54, 57, 61, 68, 2, 705, 2, 164, 3,
	2, 2, 2, 4, 182, 3, 2, 2, 2, 6, 184, 3, 2, 2, 2, 8, 195, 3, 2, 2, 2, 10,
	217, 3, 2, 2, 2, 12, 219, 3, 2, 2, 2, 14, 227, 3, 2, 2, 2, 16, 238, 3,
	2, 2, 2, 18, 240, 3, 2, 2, 2, 20, 242, 3, 2, 2, 2, 22, 244, 3, 2, 2, 2,
	24, 246, 3, 2, 2, 2, 26, 248, 3, 2, 2, 2, 28, 250, 3, 2, 2, 2, 30, 252,
	3, 2, 2, 2, 32, 254, 3, 2, 2, 2, 34, 265, 3, 2, 2, 2, 36, 269, 3, 2, 2,
	2, 38, 272, 3, 2, 2, 2, 40, 275, 3, 2, 2, 2, 42, 283, 3, 2, 2, 2, 44, 291,
	3, 2, 2, 2, 46, 296, 3, 2, 2, 2, 48, 308, 3, 2, 2, 2, 50, 317, 3, 2, 2,
	2, 52, 325, 3, 2, 2, 2, 54, 330, 3, 2, 2, 2, 56, 340, 3, 2, 2, 2, 58, 348,
	3, 2, 2, 2, 60, 353, 3, 2, 2, 2, 62, 356, 3, 2, 2, 2, 64, 363, 3, 2, 2,
	2, 66, 365, 3, 2, 2, 2, 68, 367, 3, 2, 2, 2, 70, 369, 3, 2, 2, 2, 72, 372,
	3, 2, 2, 2, 74, 395, 3, 2, 2, 2, 76, 403, 3, 2, 2, 2, 78, 407, 3, 2, 2,
	2, 80, 410, 3, 2, 2, 2, 82, 413, 3, 2, 2, 2, 84, 419, 3, 2, 2, 2, 86, 436,
	3, 2, 2, 2, 88, 438, 3, 2, 2, 2, 90, 461, 3, 2, 2, 2, 92, 471, 3, 2, 2,
	2, 94, 479, 3, 2, 2, 2, 96, 484, 3, 2, 2, 2, 98, 490, 3, 2, 2, 2, 100,
	494, 3, 2, 2, 2, 102, 501, 3, 2, 2, 2, 104, 514, 3, 2, 2, 2, 106, 528,
	3, 2, 2, 2, 108, 530, 3, 2, 2, 2, 110, 532, 3, 2, 2, 2, 112, 536, 3, 2,
	2, 2, 114, 540, 3, 2, 2, 2, 116, 547, 3, 2, 2, 2, 118, 555, 3, 2, 2, 2,
	120, 564, 3, 2, 2, 2, 122, 575, 3, 2, 2, 2, 124, 577, 3, 2, 2, 2, 126,
	579, 3, 2, 2, 2, 128, 591, 3, 2, 2, 2, 130, 601, 3, 2, 2, 2, 132, 620,
	3, 2, 2, 2, 134, 623, 3, 2, 2, 2, 136, 625, 3, 2, 2, 2, 138, 632, 3, 2,
	2, 2, 140, 642, 3, 2, 2, 2, 142, 650, 3, 2, 2, 2, 144, 652, 3, 2, 2, 2,
	146, 657, 3, 2, 2, 2, 148, 662, 3, 2, 2, 2, 150, 666, 3, 2, 2, 2, 152,
	669, 3, 2, 2, 2, 154, 671, 3, 2, 2, 2, 156, 673, 3, 2, 2, 2, 158, 675,
	3, 2, 2, 2, 160, 679, 3, 2, 2, 2, 162, 691, 3, 2, 2, 2, 164, 165, 5, 4,
	3, 2, 165, 166, 7, 2, 2, 3, 166, 3, 3, 2, 2, 2, 167, 183, 5, 6, 4, 2, 168,
	183, 5, 32, 17, 2, 169, 183, 5, 34, 18, 2, 170, 183, 5, 36, 19, 2, 171,
	183, 5, 38, 20, 2, 172, 183, 5, 40, 21, 2, 173, 183, 5, 44, 23, 2, 174,
	183, 5, 42, 22, 2, 175, 183, 5, 52, 27, 2, 176, 183, 5, 46, 24, 2, 177,
	183, 5, 48, 25, 2, 178, 183, 5, 50, 26, 2, 179, 183, 5, 54, 28, 2, 180,
	183, 5, 62, 32, 2, 181, 183, 5, 72, 37, 2, 182, 167, 3, 2, 2, 2, 182, 168,
	3, 2, 2, 2, 182, 169, 3, 2, 2, 2, 182, 170, 3, 2, 2, 2, 182, 171, 3, 2,
	2, 2, 182, 172, 3, 2, 2, 2, 182, 173, 3, 2, 2, 2, 182, 174, 3, 2, 2, 2,
	182, 175, 3, 2, 2, 2, 182, 176, 3, 2, 2, 2, 182, 177, 3, 2, 2, 2, 182,
	178, 3, 2, 2, 2, 182, 179, 3, 2, 2, 2, 182, 180, 3, 2, 2, 2, 182, 181,
	3, 2, 2, 2, 183, 5, 3, 2, 2, 2, 184, 185, 7, 3, 2, 2, 185, 186, 7, 18,
	2, 2, 186, 189, 5, 30, 16, 2, 187, 188, 7, 28, 2, 2, 188, 190, 5, 8, 5,
	2, 189, 187, 3, 2, 2, 2, 189, 190, 3, 2, 2, 2, 190, 193, 3, 2, 2, 2, 191,
	192, 7, 79, 2, 2, 192, 194, 5, 12, 7, 2, 193, 191, 3, 2, 2, 2, 193, 194,
	3, 2, 2, 2, 194, 7, 3, 2, 2, 2, 195, 200, 5, 10, 6, 2, 196, 197, 7, 79,
	2, 2, 197, 199, 5, 10, 6, 2, 198, 196, 3, 2, 2, 2, 199, 202, 3, 2, 2, 2,
	200, 198, 3, 2, 2, 2, 200, 201, 3, 2, 2, 2, 201, 9, 3, 2, 2, 2, 202, 200,
	3, 2, 2, 2, 203, 204, 7, 7, 2, 2, 204, 218, 5, 132, 67, 2, 205, 206, 7,
	9, 2, 2, 206, 218, 5, 16, 9, 2, 207, 208, 7, 10, 2, 2, 208, 218, 5, 28,
	15, 2, 209, 210, 7, 11, 2, 2, 210, 218, 5, 18, 10, 2, 211, 212, 7, 12,
	2, 2, 212, 218, 5, 20, 11, 2, 213, 214, 7, 13, 2, 2, 214, 218, 5, 22, 12,
	2, 215, 216, 7, 14, 2, 2, 216, 218, 5, 24, 13, 2, 217, 203, 3, 2, 2, 2,
	217, 205, 3, 2, 2, 2, 217, 207, 3, 2, 2, 2, 217, 209, 3, 2, 2, 2, 217,
	211, 3, 2, 2, 2, 217, 213, 3, 2, 2, 2, 217, 215, 3, 2, 2, 2, 218, 11, 3,
	2, 2, 2, 219, 224, 5, 14, 8, 2, 220, 221, 7, 79, 2, 2, 221, 223, 5, 14,
	8, 2, 222, 220, 3, 2, 2, 2, 223, 226, 3, 2, 2, 2, 224, 222, 3, 2, 2, 2,
	224, 225, 3, 2, 2, 2, 225, 13, 3, 2, 2, 2, 226, 224, 3, 2, 2, 2, 227, 228,
	7, 84, 2, 2, 228, 229, 7, 8, 2, 2, 229, 230, 5, 26, 14, 2, 230, 231, 7,
	79, 2, 2, 231, 232, 7, 11, 2, 2, 232, 233, 5, 18, 10, 2, 233, 234, 7, 79,
	2, 2, 234, 235, 7, 7, 2, 2, 235, 236, 5, 132, 67, 2, 236, 237, 7, 85, 2,
	2, 237, 15, 3, 2, 2, 2, 238, 239, 5, 146, 74, 2, 239, 17, 3, 2, 2, 2, 240,
	241, 5, 132, 67, 2, 241, 19, 3, 2, 2, 2, 242, 243, 5, 132, 67, 2, 243,
	21, 3, 2, 2, 2, 244, 245, 5, 132, 67, 2, 245, 23, 3, 2, 2, 2, 246, 247,
	5, 132, 67, 2, 247, 25, 3, 2, 2, 2, 248, 249, 5, 160, 81, 2, 249, 27, 3,
	2, 2, 2, 250, 251, 5, 146, 74, 2, 251, 29, 3, 2, 2, 2, 252, 253, 5, 160,
	81, 2, 253, 31, 3, 2, 2, 2, 254, 255, 7, 4, 2, 2, 255, 256, 7, 18, 2, 2,
	256, 259, 5, 30, 16, 2, 257, 258, 7, 28, 2, 2, 258, 260, 5, 8, 5, 2, 259,
	257, 3, 2, 2, 2, 259, 260, 3, 2, 2, 2, 260, 263, 3, 2, 2, 2, 261, 262,
	7, 79, 2, 2, 262, 264, 5, 12, 7, 2, 263, 261, 3, 2, 2, 2, 263, 264, 3,
	2, 2, 2, 264, 33, 3, 2, 2, 2, 265, 266, 7, 6, 2, 2, 266, 267, 7, 18, 2,
	2, 267, 268, 5, 30, 16, 2, 268, 35, 3, 2, 2, 2, 269, 270, 7, 17, 2, 2,
	270, 271, 7, 19, 2, 2, 271, 37, 3, 2, 2, 2, 272, 273, 7, 17, 2, 2, 273,
	274, 7, 20, 2, 2, 274, 39, 3, 2, 2, 2, 275, 276, 7, 17, 2, 2, 276, 278,
	7, 21, 2, 2, 277, 279, 5, 56, 29, 2, 278, 277, 3, 2, 2, 2, 278, 279, 3,
	2, 2, 2, 279, 281, 3, 2, 2, 2, 280, 282, 5, 150, 76, 2, 281, 280, 3, 2,
	2, 2, 281, 282, 3, 2, 2, 2, 282, 41, 3, 2, 2, 2, 283, 284, 7, 17, 2, 2,
	284, 285, 7, 24, 2, 2, 285, 286, 7, 26, 2, 2, 286, 287, 7, 31, 2, 2, 287,
	289, 5, 152, 77, 2, 288, 290, 5, 150, 76, 2, 289, 288, 3, 2, 2, 2, 289,
	290, 3, 2, 2, 2, 290, 43, 3, 2, 2, 2, 291, 292, 7, 17, 2, 2, 292, 293,
	7, 25, 2, 2, 293, 294, 7, 31, 2, 2, 294, 295, 5, 152, 77, 2, 295, 45, 3,
	2, 2, 2, 296, 297, 7, 17, 2, 2, 297, 298, 7, 24, 2, 2, 298, 299, 7, 29,
	2, 2, 299, 300, 7, 31, 2, 2, 300, 301, 5, 152, 77, 2, 301, 303, 5, 58,
	30, 2, 302, 304, 5, 60, 31, 2, 303, 302, 3, 2, 2, 2, 303, 304, 3, 2, 2,
	2, 304, 306, 3, 2, 2, 2, 305, 307, 5, 150, 76, 2, 306, 305, 3, 2, 2, 2,
	306, 307, 3, 2, 2, 2, 307, 47, 3, 2, 2, 2, 308, 309, 7, 17, 2, 2, 309,
	310, 7, 24, 2, 2, 310, 311, 7, 29, 2, 2, 311, 312, 7, 25, 2, 2, 312, 313,
	7, 31, 2, 2, 313, 314, 5, 152, 77, 2, 314, 315, 5, 58, 30, 2, 315, 316,
	5, 60, 31, 2, 316, 49, 3, 2, 2, 2, 317, 318, 7, 17, 2, 2, 318, 319, 7,
	23, 2, 2, 319, 320, 7, 26, 2, 2, 320, 321, 7, 31, 2, 2, 321, 323, 5, 152,
	77, 2, 322, 324, 5, 150, 76, 2, 323, 322, 3, 2, 2, 2, 323, 324, 3, 2, 2,
	2, 324, 51, 3, 2, 2, 2, 325, 326, 7, 17, 2, 2, 326, 328, 7, 34, 2, 2, 327,
	329, 5, 150, 76, 2, 328, 327, 3, 2, 2, 2, 328, 329, 3, 2, 2, 2, 329, 53,
	3, 2, 2, 2, 330, 331, 7, 17, 2, 2, 331, 334, 7, 56, 2, 2, 332, 333, 7,
	55, 2, 2, 333, 335, 5, 68, 35, 2, 334, 332, 3, 2, 2, 2, 334, 335, 3, 2,
	2, 2, 335, 338, 3, 2, 2, 2, 336, 337, 7, 28, 2, 2, 337, 339, 5, 70, 36,
	2, 338, 336, 3, 2, 2, 2, 338, 339, 3, 2, 2, 2, 339, 55, 3, 2, 2, 2, 340,
	341, 7, 28, 2, 2, 341, 346, 7, 22, 2, 2, 342, 343, 7, 71, 2, 2, 343, 347,
	5, 152, 77, 2, 344, 345, 7, 78, 2, 2, 345, 347, 5, 152, 77, 2, 346, 342,
	3, 2, 2, 2, 346, 344, 3, 2, 2, 2, 347, 57, 3, 2, 2, 2, 348, 349, 7, 28,
	2, 2, 349, 350, 7, 27, 2, 2, 350, 351, 7, 71, 2, 2, 351, 352, 5, 154, 78,
	2, 352, 59, 3, 2, 2, 2, 353, 354, 7, 32, 2, 2, 354, 355, 5, 86, 44, 2,
	355, 61, 3, 2, 2, 2, 356, 357, 7, 15, 2, 2, 357, 358, 7, 35, 2, 2, 358,
	361, 5, 64, 33, 2, 359, 360, 7, 16, 2, 2, 360, 362, 5, 66, 34, 2, 361,
	359, 3, 2, 2, 2, 361, 362, 3, 2, 2, 2, 362, 63, 3, 2, 2, 2, 363, 364, 7,
	92, 2, 2, 364, 65, 3, 2, 2, 2, 365, 366, 7, 92, 2, 2, 366, 67, 3, 2, 2,
	2, 367, 368, 5, 160, 81, 2, 368, 69, 3, 2, 2, 2, 369, 370, 5, 160, 81,
	2, 370, 71, 3, 2, 2, 2, 371, 373, 7, 36, 2, 2, 372, 371, 3, 2, 2, 2, 372,
	373, 3, 2, 2, 2, 373, 374, 3, 2, 2, 2, 374, 375, 7, 38, 2, 2, 375, 376,
	5, 74, 38, 2, 376, 378, 5, 80, 41, 2, 377, 379, 5, 82, 42, 2, 378, 377,
	3, 2, 2, 2, 378, 379, 3, 2, 2, 2, 379, 381, 3, 2, 2, 2, 380, 382, 5, 102,
	52, 2, 381, 380, 3, 2, 2, 2, 381, 382, 3, 2, 2, 2, 382, 384, 3, 2, 2, 2,
	383, 385, 5, 112, 57, 2, 384, 383, 3, 2, 2, 2, 384, 385, 3, 2, 2, 2, 385,
	387, 3, 2, 2, 2, 386, 388, 5, 110, 56, 2, 387, 386, 3, 2, 2, 2, 387, 388,
	3, 2, 2, 2, 388, 390, 3, 2, 2, 2, 389, 391, 5, 150, 76, 2, 390, 389, 3,
	2, 2, 2, 390, 391, 3, 2, 2, 2, 391, 393, 3, 2, 2, 2, 392, 394, 7, 37, 2,
	2, 393, 392, 3, 2, 2, 2, 393, 394, 3, 2, 2, 2, 394, 73, 3, 2, 2, 2, 395,
	400, 5, 76, 39, 2, 396, 397, 7, 79, 2, 2, 397, 399, 5, 76, 39, 2, 398,
	396, 3, 2, 2, 2, 399, 402, 3, 2, 2, 2, 400, 398, 3, 2, 2, 2, 400, 401,
	3, 2, 2, 2, 401, 75, 3, 2, 2, 2, 402, 400, 3, 2, 2, 2, 403, 405, 5, 130,
	66, 2, 404, 406, 5, 78, 40, 2, 405, 404, 3, 2, 2, 2, 405, 406, 3, 2, 2,
	2, 406, 77, 3, 2, 2, 2, 407, 408, 7, 39, 2, 2, 408, 409, 5, 160, 81, 2,
	409, 79, 3, 2, 2, 2, 410, 411, 7, 31, 2, 2, 411, 412, 5, 152, 77, 2, 412,
	81, 3, 2, 2, 2, 413, 414, 7, 32, 2, 2, 414, 415, 5, 84, 43, 2, 415, 83,
	3, 2, 2, 2, 416, 417, 8, 43, 1, 2, 417, 420, 5, 90, 46, 2, 418, 420, 5,
	94, 48, 2, 419, 416, 3, 2, 2, 2, 419, 418, 3, 2, 2, 2, 420, 426, 3, 2,
	2, 2, 421, 422, 12, 3, 2, 2, 422, 423, 7, 40, 2, 2, 423, 425, 5, 84, 43,
	4, 424, 421, 3, 2, 2, 2, 425, 428, 3, 2, 2, 2, 426, 424, 3, 2, 2, 2, 426,
	427, 3, 2, 2, 2, 427, 85, 3, 2, 2, 2, 428, 426, 3, 2, 2, 2, 429, 437, 5,
	88, 45, 2, 430, 437, 5, 90, 46, 2, 431, 434, 5, 88, 45, 2, 432, 433, 7,
	40, 2, 2, 433, 435, 5, 90, 46, 2, 434, 432, 3, 2, 2, 2, 434, 435, 3, 2,
	2, 2, 435, 437, 3, 2, 2, 2, 436, 429, 3, 2, 2, 2, 436, 430, 3, 2, 2, 2,
	436, 431, 3, 2, 2, 2, 437, 87, 3, 2, 2, 2, 438, 439, 7, 30, 2, 2, 439,
	440, 7, 71, 2, 2, 440, 441, 5, 158, 80, 2, 441, 89, 3, 2, 2, 2, 442, 443,
	8, 46, 1, 2, 443, 444, 7, 84, 2, 2, 444, 445, 5, 90, 46, 2, 445, 446, 7,
	85, 2, 2, 446, 462, 3, 2, 2, 2, 447, 448, 5, 154, 78, 2, 448, 449, 9, 2,
	2, 2, 449, 450, 5, 156, 79, 2, 450, 462, 3, 2, 2, 2, 451, 455, 5, 154,
	78, 2, 452, 456, 7, 59, 2, 2, 453, 454, 7, 49, 2, 2, 454, 456, 7, 59, 2,
	2, 455, 452, 3, 2, 2, 2, 455, 453, 3, 2, 2, 2, 456, 457, 3, 2, 2, 2, 457,
	458, 7, 84, 2, 2, 458, 459, 5, 92, 47, 2, 459, 460, 7, 85, 2, 2, 460, 462,
	3, 2, 2, 2, 461, 442, 3, 2, 2, 2, 461, 447, 3, 2, 2, 2, 461, 451, 3, 2,
	2, 2, 462, 468, 3, 2, 2, 2, 463, 464, 12, 3, 2, 2, 464, 465, 9, 3, 2, 2,
	465, 467, 5, 90, 46, 4, 466, 463, 3, 2, 2, 2, 467, 470, 3, 2, 2, 2, 468,
	466, 3, 2, 2, 2, 468, 469, 3, 2, 2, 2, 469, 91, 3, 2, 2, 2, 470, 468, 3,
	2, 2, 2, 471, 476, 5, 156, 79, 2, 472, 473, 7, 79, 2, 2, 473, 475, 5, 156,
	79, 2, 474, 472, 3, 2, 2, 2, 475, 478, 3, 2, 2, 2, 476, 474, 3, 2, 2, 2,
	476, 477, 3, 2, 2, 2, 477, 93, 3, 2, 2, 2, 478, 476, 3, 2, 2, 2, 479, 482,
	5, 96, 49, 2, 480, 481, 7, 40, 2, 2, 481, 483, 5, 96, 49, 2, 482, 480,
	3, 2, 2, 2, 482, 483, 3, 2, 2, 2, 483, 95, 3, 2, 2, 2, 484, 485, 7, 57,
	2, 2, 485, 488, 5, 128, 65, 2, 486, 489, 5, 98, 50, 2, 487, 489, 5, 160,
	81, 2, 488, 486, 3, 2, 2, 2, 488, 487, 3, 2, 2, 2, 489, 97, 3, 2, 2, 2,
	490, 492, 5, 100, 51, 2, 491, 493, 5, 132, 67, 2, 492, 491, 3, 2, 2, 2,
	492, 493, 3, 2, 2, 2, 493, 99, 3, 2, 2, 2, 494, 495, 7, 58, 2, 2, 495,
	497, 7, 84, 2, 2, 496, 498, 5, 138, 70, 2, 497, 496, 3, 2, 2, 2, 497, 498,
	3, 2, 2, 2, 498, 499, 3, 2, 2, 2, 499, 500, 7, 85, 2, 2, 500, 101, 3, 2,
	2, 2, 501, 502, 7, 52, 2, 2, 502, 503, 7, 54, 2, 2, 503, 509, 5, 104, 53,
	2, 504, 505, 7, 42, 2, 2, 505, 506, 7, 84, 2, 2, 506, 507, 5, 108, 55,
	2, 507, 508, 7, 85, 2, 2, 508, 510, 3, 2, 2, 2, 509, 504, 3, 2, 2, 2, 509,
	510, 3, 2, 2, 2, 510, 512, 3, 2, 2, 2, 511, 513, 5, 118, 60, 2, 512, 511,
	3, 2, 2, 2, 512, 513, 3, 2, 2, 2, 513, 103, 3, 2, 2, 2, 514, 519, 5, 106,
	54, 2, 515, 516, 7, 79, 2, 2, 516, 518, 5, 106, 54, 2, 517, 515, 3, 2,
	2, 2, 518, 521, 3, 2, 2, 2, 519, 517, 3, 2, 2, 2, 519, 520, 3, 2, 2, 2,
	520, 105, 3, 2, 2, 2, 521, 519, 3, 2, 2, 2, 522, 529, 5, 160, 81, 2, 523,
	524, 7, 57, 2, 2, 524, 525, 7, 84, 2, 2, 525, 526, 5, 132, 67, 2, 526,
	527, 7, 85, 2, 2, 527, 529, 3, 2, 2, 2, 528, 522, 3, 2, 2, 2, 528, 523,
	3, 2, 2, 2, 529, 107, 3, 2, 2, 2, 530, 531, 9, 4, 2, 2, 531, 109, 3, 2,
	2, 2, 532, 533, 7, 45, 2, 2, 533, 534, 7, 54, 2, 2, 534, 535, 5, 116, 59,
	2, 535, 111, 3, 2, 2, 2, 536, 537, 7, 7, 2, 2, 537, 538, 7, 54, 2, 2, 538,
	539, 5, 26, 14, 2, 539, 113, 3, 2, 2, 2, 540, 544, 5, 130, 66, 2, 541,
	543, 9, 5, 2, 2, 542, 541, 3, 2, 2, 2, 543, 546, 3, 2, 2, 2, 544, 542,
	3, 2, 2, 2, 544, 545, 3, 2, 2, 2, 545, 115, 3, 2, 2, 2, 546, 544, 3, 2,
	2, 2, 547, 552, 5, 114, 58, 2, 548, 549, 7, 79, 2, 2, 549, 551, 5, 114,
	58, 2, 550, 548, 3, 2, 2, 2, 551, 554, 3, 2, 2, 2, 552, 550, 3, 2, 2, 2,
	552, 553, 3, 2, 2, 2, 553, 117, 3, 2, 2, 2, 554, 552, 3, 2, 2, 2, 555,
	556, 7, 53, 2, 2, 556, 557, 5, 120, 61, 2, 557, 119, 3, 2, 2, 2, 558, 559,
	8, 61, 1, 2, 559, 560, 7, 84, 2, 2, 560, 561, 5, 120, 61, 2, 561, 562,
	7, 85, 2, 2, 562, 565, 3, 2, 2, 2, 563, 565, 5, 124, 63, 2, 564, 558, 3,
	2, 2, 2, 564, 563, 3, 2, 2, 2, 565, 572, 3, 2, 2, 2, 566, 567, 12, 4, 2,
	2, 567, 568, 5, 122, 62, 2, 568, 569, 5, 120, 61, 5, 569, 571, 3, 2, 2,
	2, 570, 566, 3, 2, 2, 2, 571, 574, 3, 2, 2, 2, 572, 570, 3, 2, 2, 2, 572,
	573, 3, 2, 2, 2, 573, 121, 3, 2, 2, 2, 574, 572, 3, 2, 2, 2, 575, 576,
	9, 3, 2, 2, 576, 123, 3, 2, 2, 2, 577, 578, 5, 126, 64, 2, 578, 125, 3,
	2, 2, 2, 579, 580, 5, 130, 66, 2, 580, 581, 5, 128, 65, 2, 581, 582, 5,
	130, 66, 2, 582, 127, 3, 2, 2, 2, 583, 592, 7, 71, 2, 2, 584, 592, 7, 72,
	2, 2, 585, 592, 7, 73, 2, 2, 586, 592, 7, 76, 2, 2, 587, 592, 7, 77, 2,
	2, 588, 592, 7, 74, 2, 2, 589, 592, 7, 75, 2, 2, 590, 592, 9, 6, 2, 2,
	591, 583, 3, 2, 2, 2, 591, 584, 3, 2, 2, 2, 591, 585, 3, 2, 2, 2, 591,
	586, 3, 2, 2, 2, 591, 587, 3, 2, 2, 2, 591, 588, 3, 2, 2, 2, 591, 589,
	3, 2, 2, 2, 591, 590, 3, 2, 2, 2, 592, 129, 3, 2, 2, 2, 593, 594, 8, 66,
	1, 2, 594, 595, 7, 84, 2, 2, 595, 596, 5, 130, 66, 2, 596, 597, 7, 85,
	2, 2, 597, 602, 3, 2, 2, 2, 598, 602, 5, 136, 69, 2, 599, 602, 5, 142,
	72, 2, 600, 602, 5, 132, 67, 2, 601, 593, 3, 2, 2, 2, 601, 598, 3, 2, 2,
	2, 601, 599, 3, 2, 2, 2, 601, 600, 3, 2, 2, 2, 602, 617, 3, 2, 2, 2, 603,
	604, 12, 10, 2, 2, 604, 605, 7, 89, 2, 2, 605, 616, 5, 130, 66, 11, 606,
	607, 12, 9, 2, 2, 607, 608, 7, 88, 2, 2, 608, 616, 5, 130, 66, 10, 609,
	610, 12, 8, 2, 2, 610, 611, 7, 86, 2, 2, 611, 616, 5, 130, 66, 9, 612,
	613, 12, 7, 2, 2, 613, 614, 7, 87, 2, 2, 614, 616, 5, 130, 66, 8, 615,
	603, 3, 2, 2, 2, 615, 606, 3, 2, 2, 2, 615, 609, 3, 2, 2, 2, 615, 612,
	3, 2, 2, 2, 616, 619, 3, 2, 2, 2, 617, 615, 3, 2, 2, 2, 617, 618, 3, 2,
	2, 2, 618, 131, 3, 2, 2, 2, 619, 617, 3, 2, 2, 2, 620, 621, 5, 146, 74,
	2, 621, 622, 5, 134, 68, 2, 622, 133, 3, 2, 2, 2, 623, 624, 9, 7, 2, 2,
	624, 135, 3, 2, 2, 2, 625, 626, 5, 160, 81, 2, 626, 628, 7, 84, 2, 2, 627,
	629, 5, 138, 70, 2, 628, 627, 3, 2, 2, 2, 628, 629, 3, 2, 2, 2, 629, 630,
	3, 2, 2, 2, 630, 631, 7, 85, 2, 2, 631, 137, 3, 2, 2, 2, 632, 637, 5, 140,
	71, 2, 633, 634, 7, 79, 2, 2, 634, 636, 5, 140, 71, 2, 635, 633, 3, 2,
	2, 2, 636, 639, 3, 2, 2, 2, 637, 635, 3, 2, 2, 2, 637, 638, 3, 2, 2, 2,
	638, 139, 3, 2, 2, 2, 639, 637, 3, 2, 2, 2, 640, 643, 5, 130, 66, 2, 641,
	643, 5, 90, 46, 2, 642, 640, 3, 2, 2, 2, 642, 641, 3, 2, 2, 2, 643, 141,
	3, 2, 2, 2, 644, 646, 5, 160, 81, 2, 645, 647, 5, 144, 73, 2, 646, 645,
	3, 2, 2, 2, 646, 647, 3, 2, 2, 2, 647, 651, 3, 2, 2, 2, 648, 651, 5, 148,
	75, 2, 649, 651, 5, 146, 74, 2, 650, 644, 3, 2, 2, 2, 650, 648, 3, 2, 2,
	2, 650, 649, 3, 2, 2, 2, 651, 143, 3, 2, 2, 2, 652, 653, 7, 82, 2, 2, 653,
	654, 5, 90, 46, 2, 654, 655, 7, 83, 2, 2, 655, 145, 3, 2, 2, 2, 656, 658,
	9, 8, 2, 2, 657, 656, 3, 2, 2, 2, 657, 658, 3, 2, 2, 2, 658, 659, 3, 2,
	2, 2, 659, 660, 7, 92, 2, 2, 660, 147, 3, 2, 2, 2, 661, 663, 9, 8, 2, 2,
	662, 661, 3, 2, 2, 2, 662, 663, 3, 2, 2, 2, 663, 664, 3, 2, 2, 2, 664,
	665, 7, 93, 2, 2, 665, 149, 3, 2, 2, 2, 666, 667, 7, 33, 2, 2, 667, 668,
	7, 92, 2, 2, 668, 151, 3, 2, 2, 2, 669, 670, 5, 160, 81, 2, 670, 153, 3,
	2, 2, 2, 671, 672, 5, 160, 81, 2, 672, 155, 3, 2, 2, 2, 673, 674, 5, 160,
	81, 2, 674, 157, 3, 2, 2, 2, 675, 676, 5, 160, 81, 2, 676, 159, 3, 2, 2,
	2, 677, 680, 7, 91, 2, 2, 678, 680, 5, 162, 82, 2, 679, 677, 3, 2, 2, 2,
	679, 678, 3, 2, 2, 2, 680, 688, 3, 2, 2, 2, 681, 684, 7, 69, 2, 2, 682,
	685, 7, 91, 2, 2, 683, 685, 5, 162, 82, 2, 684, 682, 3, 2, 2, 2, 684, 683,
	3, 2, 2, 2, 685, 687, 3, 2, 2, 2, 686, 681, 3, 2, 2, 2, 687, 690, 3, 2,
	2, 2, 688, 686, 3, 2, 2, 2, 688, 689, 3, 2, 2, 2, 689, 161, 3, 2, 2, 2,
	690, 688, 3, 2, 2, 2, 691, 692, 9, 9, 2, 2, 692, 163, 3, 2, 2, 2, 64, 182,
	189, 193, 200, 217, 224, 259, 263, 278, 281, 289, 303, 306, 323, 328, 334,
	338, 346, 361, 372, 378, 381, 384, 387, 390, 393, 400, 405, 419, 426, 434,
	436, 455, 461, 468, 476, 482, 488, 492, 497, 509, 512, 519, 528, 544, 552,
	564, 572, 591, 601, 615, 617, 628, 637, 642, 646, 650, 657, 662, 679, 684,
	688,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "'m'", "", "", "", "'M'", "", "'.'", "':'",
	"'='", "'<>'", "'!='", "'>'", "'>='", "'<'", "'<='", "'=~'", "','", "'{'",
	"'}'", "'['", "']'", "'('", "')'", "'+'", "'-'", "'/'", "'*'", "'%'",
}
var symbolicNames = []string{
	"", "T_CREATE", "T_UPDATE", "T_SET", "T_DROP", "T_INTERVAL", "T_INTERVAL_NAME",
	"T_SHARD", "T_REPLICATION", "T_TTL", "T_META_TTL", "T_PAST_TTL", "T_FUTURE_TTL",
	"T_KILL", "T_ON", "T_SHOW", "T_DATASBAE", "T_DATASBAES", "T_NODE", "T_MEASUREMENTS",
	"T_MEASUREMENT", "T_FIELD", "T_TAG", "T_INFO", "T_KEYS", "T_KEY", "T_WITH",
	"T_VALUES", "T_VALUE", "T_FROM", "T_WHERE", "T_LIMIT", "T_QUERIES", "T_QUERY",
	"T_EXPLAIN", "T_WITH_VALUE", "T_SELECT", "T_AS", "T_AND", "T_OR", "T_FILL",
	"T_NULL", "T_PREVIOUS", "T_ORDER", "T_ASC", "T_DESC", "T_LIKE", "T_NOT",
	"T_BETWEEN", "T_IS", "T_GROUP", "T_HAVING", "T_BY", "T_FOR", "T_STATS",
	"T_TIME", "T_NOW", "T_IN", "T_LOG", "T_PROFILE", "T_SECOND", "T_MINUTE",
	"T_HOUR", "T_DAY", "T_WEEK", "T_MONTH", "T_YEAR", "T_DOT", "T_COLON", "T_EQUAL",
	"T_NOTEQUAL", "T_NOTEQUAL2", "T_GREATER", "T_GREATEREQUAL", "T_LESS", "T_LESSEQUAL",
	"T_REGEXP", "T_COMMA", "T_OPEN_B", "T_CLOSE_B", "T_OPEN_SB", "T_CLOSE_SB",
	"T_OPEN_P", "T_CLOSE_P", "T_ADD", "T_SUB", "T_DIV", "T_MUL", "T_MOD", "L_ID",
	"L_INT", "L_DEC", "WS",
}

var ruleNames = []string{
	"statement", "statementList", "createDatabaseStmt", "withClauseList", "withClause",
	"intervalDefineList", "intervalDefine", "shardNum", "ttlVal", "metattlVal",
	"pastVal", "futureVal", "intervalNameVal", "replicaFactor", "databaseName",
	"updateDatabaseStmt", "dropDatabaseStmt", "showDatabasesStmt", "showNodeStmt",
	"showMeasurementsStmt", "showTagKeysStmt", "showInfoStmt", "showTagValuesStmt",
	"showTagValuesInfoStmt", "showFieldKeysStmt", "showQueriesStmt", "showStatsStmt",
	"withMeasurementClause", "withTagClause", "whereTagCascade", "killQueryStmt",
	"queryId", "serverId", "module", "component", "queryStmt", "fields", "field",
	"alias", "fromClause", "whereClause", "clauseBooleanExpr", "tagCascadeExpr",
	"tagEqualExpr", "tagBooleanExpr", "tagValueList", "timeExpr", "timeBooleanExpr",
	"nowExpr", "nowFunc", "groupByClause", "dimensions", "dimension", "fillOption",
	"orderByClause", "intervalByClause", "sortField", "sortFields", "havingClause",
	"boolExpr", "boolExprLogicalOp", "boolExprAtom", "boolExprBinary", "boolExprBinaryOperator",
	"expr", "durationLit", "intervalItem", "exprFunc", "exprFuncParams", "funcParam",
	"exprAtom", "identFilter", "intNumber", "decNumber", "limitClause", "metricName",
	"tagKey", "tagValue", "tagValuePattern", "ident", "nonReservedWords",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type SQLParser struct {
	*antlr.BaseParser
}

func NewSQLParser(input antlr.TokenStream) *SQLParser {
	this := new(SQLParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "SQL.g4"

	return this
}

// SQLParser tokens.
const (
	SQLParserEOF             = antlr.TokenEOF
	SQLParserT_CREATE        = 1
	SQLParserT_UPDATE        = 2
	SQLParserT_SET           = 3
	SQLParserT_DROP          = 4
	SQLParserT_INTERVAL      = 5
	SQLParserT_INTERVAL_NAME = 6
	SQLParserT_SHARD         = 7
	SQLParserT_REPLICATION   = 8
	SQLParserT_TTL           = 9
	SQLParserT_META_TTL      = 10
	SQLParserT_PAST_TTL      = 11
	SQLParserT_FUTURE_TTL    = 12
	SQLParserT_KILL          = 13
	SQLParserT_ON            = 14
	SQLParserT_SHOW          = 15
	SQLParserT_DATASBAE      = 16
	SQLParserT_DATASBAES     = 17
	SQLParserT_NODE          = 18
	SQLParserT_MEASUREMENTS  = 19
	SQLParserT_MEASUREMENT   = 20
	SQLParserT_FIELD         = 21
	SQLParserT_TAG           = 22
	SQLParserT_INFO          = 23
	SQLParserT_KEYS          = 24
	SQLParserT_KEY           = 25
	SQLParserT_WITH          = 26
	SQLParserT_VALUES        = 27
	SQLParserT_VALUE         = 28
	SQLParserT_FROM          = 29
	SQLParserT_WHERE         = 30
	SQLParserT_LIMIT         = 31
	SQLParserT_QUERIES       = 32
	SQLParserT_QUERY         = 33
	SQLParserT_EXPLAIN       = 34
	SQLParserT_WITH_VALUE    = 35
	SQLParserT_SELECT        = 36
	SQLParserT_AS            = 37
	SQLParserT_AND           = 38
	SQLParserT_OR            = 39
	SQLParserT_FILL          = 40
	SQLParserT_NULL          = 41
	SQLParserT_PREVIOUS      = 42
	SQLParserT_ORDER         = 43
	SQLParserT_ASC           = 44
	SQLParserT_DESC          = 45
	SQLParserT_LIKE          = 46
	SQLParserT_NOT           = 47
	SQLParserT_BETWEEN       = 48
	SQLParserT_IS            = 49
	SQLParserT_GROUP         = 50
	SQLParserT_HAVING        = 51
	SQLParserT_BY            = 52
	SQLParserT_FOR           = 53
	SQLParserT_STATS         = 54
	SQLParserT_TIME          = 55
	SQLParserT_NOW           = 56
	SQLParserT_IN            = 57
	SQLParserT_LOG           = 58
	SQLParserT_PROFILE       = 59
	SQLParserT_SECOND        = 60
	SQLParserT_MINUTE        = 61
	SQLParserT_HOUR          = 62
	SQLParserT_DAY           = 63
	SQLParserT_WEEK          = 64
	SQLParserT_MONTH         = 65
	SQLParserT_YEAR          = 66
	SQLParserT_DOT           = 67
	SQLParserT_COLON         = 68
	SQLParserT_EQUAL         = 69
	SQLParserT_NOTEQUAL      = 70
	SQLParserT_NOTEQUAL2     = 71
	SQLParserT_GREATER       = 72
	SQLParserT_GREATEREQUAL  = 73
	SQLParserT_LESS          = 74
	SQLParserT_LESSEQUAL     = 75
	SQLParserT_REGEXP        = 76
	SQLParserT_COMMA         = 77
	SQLParserT_OPEN_B        = 78
	SQLParserT_CLOSE_B       = 79
	SQLParserT_OPEN_SB       = 80
	SQLParserT_CLOSE_SB      = 81
	SQLParserT_OPEN_P        = 82
	SQLParserT_CLOSE_P       = 83
	SQLParserT_ADD           = 84
	SQLParserT_SUB           = 85
	SQLParserT_DIV           = 86
	SQLParserT_MUL           = 87
	SQLParserT_MOD           = 88
	SQLParserL_ID            = 89
	SQLParserL_INT           = 90
	SQLParserL_DEC           = 91
	SQLParserWS              = 92
)

// SQLParser rules.
const (
	SQLParserRULE_statement              = 0
	SQLParserRULE_statementList          = 1
	SQLParserRULE_createDatabaseStmt     = 2
	SQLParserRULE_withClauseList         = 3
	SQLParserRULE_withClause             = 4
	SQLParserRULE_intervalDefineList     = 5
	SQLParserRULE_intervalDefine         = 6
	SQLParserRULE_shardNum               = 7
	SQLParserRULE_ttlVal                 = 8
	SQLParserRULE_metattlVal             = 9
	SQLParserRULE_pastVal                = 10
	SQLParserRULE_futureVal              = 11
	SQLParserRULE_intervalNameVal        = 12
	SQLParserRULE_replicaFactor          = 13
	SQLParserRULE_databaseName           = 14
	SQLParserRULE_updateDatabaseStmt     = 15
	SQLParserRULE_dropDatabaseStmt       = 16
	SQLParserRULE_showDatabasesStmt      = 17
	SQLParserRULE_showNodeStmt           = 18
	SQLParserRULE_showMeasurementsStmt   = 19
	SQLParserRULE_showTagKeysStmt        = 20
	SQLParserRULE_showInfoStmt           = 21
	SQLParserRULE_showTagValuesStmt      = 22
	SQLParserRULE_showTagValuesInfoStmt  = 23
	SQLParserRULE_showFieldKeysStmt      = 24
	SQLParserRULE_showQueriesStmt        = 25
	SQLParserRULE_showStatsStmt          = 26
	SQLParserRULE_withMeasurementClause  = 27
	SQLParserRULE_withTagClause          = 28
	SQLParserRULE_whereTagCascade        = 29
	SQLParserRULE_killQueryStmt          = 30
	SQLParserRULE_queryId                = 31
	SQLParserRULE_serverId               = 32
	SQLParserRULE_module                 = 33
	SQLParserRULE_component              = 34
	SQLParserRULE_queryStmt              = 35
	SQLParserRULE_fields                 = 36
	SQLParserRULE_field                  = 37
	SQLParserRULE_alias                  = 38
	SQLParserRULE_fromClause             = 39
	SQLParserRULE_whereClause            = 40
	SQLParserRULE_clauseBooleanExpr      = 41
	SQLParserRULE_tagCascadeExpr         = 42
	SQLParserRULE_tagEqualExpr           = 43
	SQLParserRULE_tagBooleanExpr         = 44
	SQLParserRULE_tagValueList           = 45
	SQLParserRULE_timeExpr               = 46
	SQLParserRULE_timeBooleanExpr        = 47
	SQLParserRULE_nowExpr                = 48
	SQLParserRULE_nowFunc                = 49
	SQLParserRULE_groupByClause          = 50
	SQLParserRULE_dimensions             = 51
	SQLParserRULE_dimension              = 52
	SQLParserRULE_fillOption             = 53
	SQLParserRULE_orderByClause          = 54
	SQLParserRULE_intervalByClause       = 55
	SQLParserRULE_sortField              = 56
	SQLParserRULE_sortFields             = 57
	SQLParserRULE_havingClause           = 58
	SQLParserRULE_boolExpr               = 59
	SQLParserRULE_boolExprLogicalOp      = 60
	SQLParserRULE_boolExprAtom           = 61
	SQLParserRULE_boolExprBinary         = 62
	SQLParserRULE_boolExprBinaryOperator = 63
	SQLParserRULE_expr                   = 64
	SQLParserRULE_durationLit            = 65
	SQLParserRULE_intervalItem           = 66
	SQLParserRULE_exprFunc               = 67
	SQLParserRULE_exprFuncParams         = 68
	SQLParserRULE_funcParam              = 69
	SQLParserRULE_exprAtom               = 70
	SQLParserRULE_identFilter            = 71
	SQLParserRULE_intNumber              = 72
	SQLParserRULE_decNumber              = 73
	SQLParserRULE_limitClause            = 74
	SQLParserRULE_metricName             = 75
	SQLParserRULE_tagKey                 = 76
	SQLParserRULE_tagValue               = 77
	SQLParserRULE_tagValuePattern        = 78
	SQLParserRULE_ident                  = 79
	SQLParserRULE_nonReservedWords       = 80
)

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_statement
	return p
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) StatementList() IStatementListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStatementListContext)
}

func (s *StatementContext) EOF() antlr.TerminalNode {
	return s.GetToken(SQLParserEOF, 0)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (s *StatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Statement() (localctx IStatementContext) {
	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, SQLParserRULE_statement)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(162)
		p.StatementList()
	}
	{
		p.SetState(163)
		p.Match(SQLParserEOF)
	}

	return localctx
}

// IStatementListContext is an interface to support dynamic dispatch.
type IStatementListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatementListContext differentiates from other interfaces.
	IsStatementListContext()
}

type StatementListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementListContext() *StatementListContext {
	var p = new(StatementListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_statementList
	return p
}

func (*StatementListContext) IsStatementListContext() {}

func NewStatementListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementListContext {
	var p = new(StatementListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_statementList

	return p
}

func (s *StatementListContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementListContext) CreateDatabaseStmt() ICreateDatabaseStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICreateDatabaseStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICreateDatabaseStmtContext)
}

func (s *StatementListContext) UpdateDatabaseStmt() IUpdateDatabaseStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUpdateDatabaseStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUpdateDatabaseStmtContext)
}

func (s *StatementListContext) DropDatabaseStmt() IDropDatabaseStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDropDatabaseStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDropDatabaseStmtContext)
}

func (s *StatementListContext) ShowDatabasesStmt() IShowDatabasesStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowDatabasesStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowDatabasesStmtContext)
}

func (s *StatementListContext) ShowNodeStmt() IShowNodeStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowNodeStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowNodeStmtContext)
}

func (s *StatementListContext) ShowMeasurementsStmt() IShowMeasurementsStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowMeasurementsStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowMeasurementsStmtContext)
}

func (s *StatementListContext) ShowInfoStmt() IShowInfoStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowInfoStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowInfoStmtContext)
}

func (s *StatementListContext) ShowTagKeysStmt() IShowTagKeysStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowTagKeysStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowTagKeysStmtContext)
}

func (s *StatementListContext) ShowQueriesStmt() IShowQueriesStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowQueriesStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowQueriesStmtContext)
}

func (s *StatementListContext) ShowTagValuesStmt() IShowTagValuesStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowTagValuesStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowTagValuesStmtContext)
}

func (s *StatementListContext) ShowTagValuesInfoStmt() IShowTagValuesInfoStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowTagValuesInfoStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowTagValuesInfoStmtContext)
}

func (s *StatementListContext) ShowFieldKeysStmt() IShowFieldKeysStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowFieldKeysStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowFieldKeysStmtContext)
}

func (s *StatementListContext) ShowStatsStmt() IShowStatsStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowStatsStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowStatsStmtContext)
}

func (s *StatementListContext) KillQueryStmt() IKillQueryStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IKillQueryStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IKillQueryStmtContext)
}

func (s *StatementListContext) QueryStmt() IQueryStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQueryStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IQueryStmtContext)
}

func (s *StatementListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterStatementList(s)
	}
}

func (s *StatementListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitStatementList(s)
	}
}

func (s *StatementListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitStatementList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) StatementList() (localctx IStatementListContext) {
	localctx = NewStatementListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, SQLParserRULE_statementList)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(180)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(165)
			p.CreateDatabaseStmt()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(166)
			p.UpdateDatabaseStmt()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(167)
			p.DropDatabaseStmt()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(168)
			p.ShowDatabasesStmt()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(169)
			p.ShowNodeStmt()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(170)
			p.ShowMeasurementsStmt()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(171)
			p.ShowInfoStmt()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(172)
			p.ShowTagKeysStmt()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(173)
			p.ShowQueriesStmt()
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(174)
			p.ShowTagValuesStmt()
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(175)
			p.ShowTagValuesInfoStmt()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(176)
			p.ShowFieldKeysStmt()
		}

	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(177)
			p.ShowStatsStmt()
		}

	case 14:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(178)
			p.KillQueryStmt()
		}

	case 15:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(179)
			p.QueryStmt()
		}

	}

	return localctx
}

// ICreateDatabaseStmtContext is an interface to support dynamic dispatch.
type ICreateDatabaseStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCreateDatabaseStmtContext differentiates from other interfaces.
	IsCreateDatabaseStmtContext()
}

type CreateDatabaseStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateDatabaseStmtContext() *CreateDatabaseStmtContext {
	var p = new(CreateDatabaseStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_createDatabaseStmt
	return p
}

func (*CreateDatabaseStmtContext) IsCreateDatabaseStmtContext() {}

func NewCreateDatabaseStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateDatabaseStmtContext {
	var p = new(CreateDatabaseStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_createDatabaseStmt

	return p
}

func (s *CreateDatabaseStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateDatabaseStmtContext) T_CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CREATE, 0)
}

func (s *CreateDatabaseStmtContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *CreateDatabaseStmtContext) DatabaseName() IDatabaseNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDatabaseNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDatabaseNameContext)
}

func (s *CreateDatabaseStmtContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *CreateDatabaseStmtContext) WithClauseList() IWithClauseListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithClauseListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWithClauseListContext)
}

func (s *CreateDatabaseStmtContext) T_COMMA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, 0)
}

func (s *CreateDatabaseStmtContext) IntervalDefineList() IIntervalDefineListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntervalDefineListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntervalDefineListContext)
}

func (s *CreateDatabaseStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateDatabaseStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CreateDatabaseStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterCreateDatabaseStmt(s)
	}
}

func (s *CreateDatabaseStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitCreateDatabaseStmt(s)
	}
}

func (s *CreateDatabaseStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitCreateDatabaseStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) CreateDatabaseStmt() (localctx ICreateDatabaseStmtContext) {
	localctx = NewCreateDatabaseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SQLParserRULE_createDatabaseStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(182)
		p.Match(SQLParserT_CREATE)
	}
	{
		p.SetState(183)
		p.Match(SQLParserT_DATASBAE)
	}
	{
		p.SetState(184)
		p.DatabaseName()
	}
	p.SetState(187)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_WITH {
		{
			p.SetState(185)
			p.Match(SQLParserT_WITH)
		}
		{
			p.SetState(186)
			p.WithClauseList()
		}

	}
	p.SetState(191)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_COMMA {
		{
			p.SetState(189)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(190)
			p.IntervalDefineList()
		}

	}

	return localctx
}

// IWithClauseListContext is an interface to support dynamic dispatch.
type IWithClauseListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWithClauseListContext differentiates from other interfaces.
	IsWithClauseListContext()
}

type WithClauseListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithClauseListContext() *WithClauseListContext {
	var p = new(WithClauseListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_withClauseList
	return p
}

func (*WithClauseListContext) IsWithClauseListContext() {}

func NewWithClauseListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithClauseListContext {
	var p = new(WithClauseListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_withClauseList

	return p
}

func (s *WithClauseListContext) GetParser() antlr.Parser { return s.parser }

func (s *WithClauseListContext) AllWithClause() []IWithClauseContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWithClauseContext)(nil)).Elem())
	var tst = make([]IWithClauseContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWithClauseContext)
		}
	}

	return tst
}

func (s *WithClauseListContext) WithClause(i int) IWithClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithClauseContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWithClauseContext)
}

func (s *WithClauseListContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *WithClauseListContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *WithClauseListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithClauseListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WithClauseListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWithClauseList(s)
	}
}

func (s *WithClauseListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWithClauseList(s)
	}
}

func (s *WithClauseListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWithClauseList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) WithClauseList() (localctx IWithClauseListContext) {
	localctx = NewWithClauseListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, SQLParserRULE_withClauseList)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(193)
		p.WithClause()
	}
	p.SetState(198)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(194)
				p.Match(SQLParserT_COMMA)
			}
			{
				p.SetState(195)
				p.WithClause()
			}

		}
		p.SetState(200)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}

	return localctx
}

// IWithClauseContext is an interface to support dynamic dispatch.
type IWithClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWithClauseContext differentiates from other interfaces.
	IsWithClauseContext()
}

type WithClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithClauseContext() *WithClauseContext {
	var p = new(WithClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_withClause
	return p
}

func (*WithClauseContext) IsWithClauseContext() {}

func NewWithClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithClauseContext {
	var p = new(WithClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_withClause

	return p
}

func (s *WithClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *WithClauseContext) T_INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL, 0)
}

func (s *WithClauseContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *WithClauseContext) T_SHARD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHARD, 0)
}

func (s *WithClauseContext) ShardNum() IShardNumContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShardNumContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShardNumContext)
}

func (s *WithClauseContext) T_REPLICATION() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REPLICATION, 0)
}

func (s *WithClauseContext) ReplicaFactor() IReplicaFactorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IReplicaFactorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IReplicaFactorContext)
}

func (s *WithClauseContext) T_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TTL, 0)
}

func (s *WithClauseContext) TtlVal() ITtlValContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITtlValContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITtlValContext)
}

func (s *WithClauseContext) T_META_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_META_TTL, 0)
}

func (s *WithClauseContext) MetattlVal() IMetattlValContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetattlValContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetattlValContext)
}

func (s *WithClauseContext) T_PAST_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PAST_TTL, 0)
}

func (s *WithClauseContext) PastVal() IPastValContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPastValContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPastValContext)
}

func (s *WithClauseContext) T_FUTURE_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FUTURE_TTL, 0)
}

func (s *WithClauseContext) FutureVal() IFutureValContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFutureValContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFutureValContext)
}

func (s *WithClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WithClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWithClause(s)
	}
}

func (s *WithClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWithClause(s)
	}
}

func (s *WithClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWithClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) WithClause() (localctx IWithClauseContext) {
	localctx = NewWithClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SQLParserRULE_withClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(215)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_INTERVAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(201)
			p.Match(SQLParserT_INTERVAL)
		}
		{
			p.SetState(202)
			p.DurationLit()
		}

	case SQLParserT_SHARD:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(203)
			p.Match(SQLParserT_SHARD)
		}
		{
			p.SetState(204)
			p.ShardNum()
		}

	case SQLParserT_REPLICATION:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(205)
			p.Match(SQLParserT_REPLICATION)
		}
		{
			p.SetState(206)
			p.ReplicaFactor()
		}

	case SQLParserT_TTL:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(207)
			p.Match(SQLParserT_TTL)
		}
		{
			p.SetState(208)
			p.TtlVal()
		}

	case SQLParserT_META_TTL:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(209)
			p.Match(SQLParserT_META_TTL)
		}
		{
			p.SetState(210)
			p.MetattlVal()
		}

	case SQLParserT_PAST_TTL:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(211)
			p.Match(SQLParserT_PAST_TTL)
		}
		{
			p.SetState(212)
			p.PastVal()
		}

	case SQLParserT_FUTURE_TTL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(213)
			p.Match(SQLParserT_FUTURE_TTL)
		}
		{
			p.SetState(214)
			p.FutureVal()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IIntervalDefineListContext is an interface to support dynamic dispatch.
type IIntervalDefineListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntervalDefineListContext differentiates from other interfaces.
	IsIntervalDefineListContext()
}

type IntervalDefineListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntervalDefineListContext() *IntervalDefineListContext {
	var p = new(IntervalDefineListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_intervalDefineList
	return p
}

func (*IntervalDefineListContext) IsIntervalDefineListContext() {}

func NewIntervalDefineListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntervalDefineListContext {
	var p = new(IntervalDefineListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_intervalDefineList

	return p
}

func (s *IntervalDefineListContext) GetParser() antlr.Parser { return s.parser }

func (s *IntervalDefineListContext) AllIntervalDefine() []IIntervalDefineContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IIntervalDefineContext)(nil)).Elem())
	var tst = make([]IIntervalDefineContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IIntervalDefineContext)
		}
	}

	return tst
}

func (s *IntervalDefineListContext) IntervalDefine(i int) IIntervalDefineContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntervalDefineContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IIntervalDefineContext)
}

func (s *IntervalDefineListContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *IntervalDefineListContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *IntervalDefineListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntervalDefineListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntervalDefineListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIntervalDefineList(s)
	}
}

func (s *IntervalDefineListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIntervalDefineList(s)
	}
}

func (s *IntervalDefineListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIntervalDefineList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) IntervalDefineList() (localctx IIntervalDefineListContext) {
	localctx = NewIntervalDefineListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SQLParserRULE_intervalDefineList)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(217)
		p.IntervalDefine()
	}
	p.SetState(222)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserT_COMMA {
		{
			p.SetState(218)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(219)
			p.IntervalDefine()
		}

		p.SetState(224)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IIntervalDefineContext is an interface to support dynamic dispatch.
type IIntervalDefineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntervalDefineContext differentiates from other interfaces.
	IsIntervalDefineContext()
}

type IntervalDefineContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntervalDefineContext() *IntervalDefineContext {
	var p = new(IntervalDefineContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_intervalDefine
	return p
}

func (*IntervalDefineContext) IsIntervalDefineContext() {}

func NewIntervalDefineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntervalDefineContext {
	var p = new(IntervalDefineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_intervalDefine

	return p
}

func (s *IntervalDefineContext) GetParser() antlr.Parser { return s.parser }

func (s *IntervalDefineContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *IntervalDefineContext) T_INTERVAL_NAME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL_NAME, 0)
}

func (s *IntervalDefineContext) IntervalNameVal() IIntervalNameValContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntervalNameValContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntervalNameValContext)
}

func (s *IntervalDefineContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *IntervalDefineContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *IntervalDefineContext) T_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TTL, 0)
}

func (s *IntervalDefineContext) TtlVal() ITtlValContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITtlValContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITtlValContext)
}

func (s *IntervalDefineContext) T_INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL, 0)
}

func (s *IntervalDefineContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *IntervalDefineContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *IntervalDefineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntervalDefineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntervalDefineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIntervalDefine(s)
	}
}

func (s *IntervalDefineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIntervalDefine(s)
	}
}

func (s *IntervalDefineContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIntervalDefine(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) IntervalDefine() (localctx IIntervalDefineContext) {
	localctx = NewIntervalDefineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SQLParserRULE_intervalDefine)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(225)
		p.Match(SQLParserT_OPEN_P)
	}
	{
		p.SetState(226)
		p.Match(SQLParserT_INTERVAL_NAME)
	}
	{
		p.SetState(227)
		p.IntervalNameVal()
	}
	{
		p.SetState(228)
		p.Match(SQLParserT_COMMA)
	}
	{
		p.SetState(229)
		p.Match(SQLParserT_TTL)
	}
	{
		p.SetState(230)
		p.TtlVal()
	}
	{
		p.SetState(231)
		p.Match(SQLParserT_COMMA)
	}
	{
		p.SetState(232)
		p.Match(SQLParserT_INTERVAL)
	}
	{
		p.SetState(233)
		p.DurationLit()
	}
	{
		p.SetState(234)
		p.Match(SQLParserT_CLOSE_P)
	}

	return localctx
}

// IShardNumContext is an interface to support dynamic dispatch.
type IShardNumContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShardNumContext differentiates from other interfaces.
	IsShardNumContext()
}

type ShardNumContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShardNumContext() *ShardNumContext {
	var p = new(ShardNumContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_shardNum
	return p
}

func (*ShardNumContext) IsShardNumContext() {}

func NewShardNumContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShardNumContext {
	var p = new(ShardNumContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_shardNum

	return p
}

func (s *ShardNumContext) GetParser() antlr.Parser { return s.parser }

func (s *ShardNumContext) IntNumber() IIntNumberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntNumberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntNumberContext)
}

func (s *ShardNumContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShardNumContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShardNumContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShardNum(s)
	}
}

func (s *ShardNumContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShardNum(s)
	}
}

func (s *ShardNumContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShardNum(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShardNum() (localctx IShardNumContext) {
	localctx = NewShardNumContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SQLParserRULE_shardNum)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(236)
		p.IntNumber()
	}

	return localctx
}

// ITtlValContext is an interface to support dynamic dispatch.
type ITtlValContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTtlValContext differentiates from other interfaces.
	IsTtlValContext()
}

type TtlValContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTtlValContext() *TtlValContext {
	var p = new(TtlValContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_ttlVal
	return p
}

func (*TtlValContext) IsTtlValContext() {}

func NewTtlValContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TtlValContext {
	var p = new(TtlValContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_ttlVal

	return p
}

func (s *TtlValContext) GetParser() antlr.Parser { return s.parser }

func (s *TtlValContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *TtlValContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TtlValContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TtlValContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTtlVal(s)
	}
}

func (s *TtlValContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTtlVal(s)
	}
}

func (s *TtlValContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTtlVal(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TtlVal() (localctx ITtlValContext) {
	localctx = NewTtlValContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SQLParserRULE_ttlVal)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(238)
		p.DurationLit()
	}

	return localctx
}

// IMetattlValContext is an interface to support dynamic dispatch.
type IMetattlValContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMetattlValContext differentiates from other interfaces.
	IsMetattlValContext()
}

type MetattlValContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetattlValContext() *MetattlValContext {
	var p = new(MetattlValContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_metattlVal
	return p
}

func (*MetattlValContext) IsMetattlValContext() {}

func NewMetattlValContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MetattlValContext {
	var p = new(MetattlValContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_metattlVal

	return p
}

func (s *MetattlValContext) GetParser() antlr.Parser { return s.parser }

func (s *MetattlValContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *MetattlValContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MetattlValContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MetattlValContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterMetattlVal(s)
	}
}

func (s *MetattlValContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitMetattlVal(s)
	}
}

func (s *MetattlValContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitMetattlVal(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) MetattlVal() (localctx IMetattlValContext) {
	localctx = NewMetattlValContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SQLParserRULE_metattlVal)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(240)
		p.DurationLit()
	}

	return localctx
}

// IPastValContext is an interface to support dynamic dispatch.
type IPastValContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPastValContext differentiates from other interfaces.
	IsPastValContext()
}

type PastValContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPastValContext() *PastValContext {
	var p = new(PastValContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_pastVal
	return p
}

func (*PastValContext) IsPastValContext() {}

func NewPastValContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PastValContext {
	var p = new(PastValContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_pastVal

	return p
}

func (s *PastValContext) GetParser() antlr.Parser { return s.parser }

func (s *PastValContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *PastValContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PastValContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PastValContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterPastVal(s)
	}
}

func (s *PastValContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitPastVal(s)
	}
}

func (s *PastValContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitPastVal(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) PastVal() (localctx IPastValContext) {
	localctx = NewPastValContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SQLParserRULE_pastVal)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(242)
		p.DurationLit()
	}

	return localctx
}

// IFutureValContext is an interface to support dynamic dispatch.
type IFutureValContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFutureValContext differentiates from other interfaces.
	IsFutureValContext()
}

type FutureValContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFutureValContext() *FutureValContext {
	var p = new(FutureValContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_futureVal
	return p
}

func (*FutureValContext) IsFutureValContext() {}

func NewFutureValContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FutureValContext {
	var p = new(FutureValContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_futureVal

	return p
}

func (s *FutureValContext) GetParser() antlr.Parser { return s.parser }

func (s *FutureValContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *FutureValContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FutureValContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FutureValContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFutureVal(s)
	}
}

func (s *FutureValContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFutureVal(s)
	}
}

func (s *FutureValContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFutureVal(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) FutureVal() (localctx IFutureValContext) {
	localctx = NewFutureValContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SQLParserRULE_futureVal)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(244)
		p.DurationLit()
	}

	return localctx
}

// IIntervalNameValContext is an interface to support dynamic dispatch.
type IIntervalNameValContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntervalNameValContext differentiates from other interfaces.
	IsIntervalNameValContext()
}

type IntervalNameValContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntervalNameValContext() *IntervalNameValContext {
	var p = new(IntervalNameValContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_intervalNameVal
	return p
}

func (*IntervalNameValContext) IsIntervalNameValContext() {}

func NewIntervalNameValContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntervalNameValContext {
	var p = new(IntervalNameValContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_intervalNameVal

	return p
}

func (s *IntervalNameValContext) GetParser() antlr.Parser { return s.parser }

func (s *IntervalNameValContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *IntervalNameValContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntervalNameValContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntervalNameValContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIntervalNameVal(s)
	}
}

func (s *IntervalNameValContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIntervalNameVal(s)
	}
}

func (s *IntervalNameValContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIntervalNameVal(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) IntervalNameVal() (localctx IIntervalNameValContext) {
	localctx = NewIntervalNameValContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, SQLParserRULE_intervalNameVal)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(246)
		p.Ident()
	}

	return localctx
}

// IReplicaFactorContext is an interface to support dynamic dispatch.
type IReplicaFactorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsReplicaFactorContext differentiates from other interfaces.
	IsReplicaFactorContext()
}

type ReplicaFactorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReplicaFactorContext() *ReplicaFactorContext {
	var p = new(ReplicaFactorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_replicaFactor
	return p
}

func (*ReplicaFactorContext) IsReplicaFactorContext() {}

func NewReplicaFactorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReplicaFactorContext {
	var p = new(ReplicaFactorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_replicaFactor

	return p
}

func (s *ReplicaFactorContext) GetParser() antlr.Parser { return s.parser }

func (s *ReplicaFactorContext) IntNumber() IIntNumberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntNumberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntNumberContext)
}

func (s *ReplicaFactorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReplicaFactorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReplicaFactorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterReplicaFactor(s)
	}
}

func (s *ReplicaFactorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitReplicaFactor(s)
	}
}

func (s *ReplicaFactorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitReplicaFactor(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ReplicaFactor() (localctx IReplicaFactorContext) {
	localctx = NewReplicaFactorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, SQLParserRULE_replicaFactor)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(248)
		p.IntNumber()
	}

	return localctx
}

// IDatabaseNameContext is an interface to support dynamic dispatch.
type IDatabaseNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDatabaseNameContext differentiates from other interfaces.
	IsDatabaseNameContext()
}

type DatabaseNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDatabaseNameContext() *DatabaseNameContext {
	var p = new(DatabaseNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_databaseName
	return p
}

func (*DatabaseNameContext) IsDatabaseNameContext() {}

func NewDatabaseNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DatabaseNameContext {
	var p = new(DatabaseNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_databaseName

	return p
}

func (s *DatabaseNameContext) GetParser() antlr.Parser { return s.parser }

func (s *DatabaseNameContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *DatabaseNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DatabaseNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DatabaseNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDatabaseName(s)
	}
}

func (s *DatabaseNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDatabaseName(s)
	}
}

func (s *DatabaseNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDatabaseName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) DatabaseName() (localctx IDatabaseNameContext) {
	localctx = NewDatabaseNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, SQLParserRULE_databaseName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(250)
		p.Ident()
	}

	return localctx
}

// IUpdateDatabaseStmtContext is an interface to support dynamic dispatch.
type IUpdateDatabaseStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUpdateDatabaseStmtContext differentiates from other interfaces.
	IsUpdateDatabaseStmtContext()
}

type UpdateDatabaseStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUpdateDatabaseStmtContext() *UpdateDatabaseStmtContext {
	var p = new(UpdateDatabaseStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_updateDatabaseStmt
	return p
}

func (*UpdateDatabaseStmtContext) IsUpdateDatabaseStmtContext() {}

func NewUpdateDatabaseStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UpdateDatabaseStmtContext {
	var p = new(UpdateDatabaseStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_updateDatabaseStmt

	return p
}

func (s *UpdateDatabaseStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *UpdateDatabaseStmtContext) T_UPDATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_UPDATE, 0)
}

func (s *UpdateDatabaseStmtContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *UpdateDatabaseStmtContext) DatabaseName() IDatabaseNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDatabaseNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDatabaseNameContext)
}

func (s *UpdateDatabaseStmtContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *UpdateDatabaseStmtContext) WithClauseList() IWithClauseListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithClauseListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWithClauseListContext)
}

func (s *UpdateDatabaseStmtContext) T_COMMA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, 0)
}

func (s *UpdateDatabaseStmtContext) IntervalDefineList() IIntervalDefineListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntervalDefineListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntervalDefineListContext)
}

func (s *UpdateDatabaseStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UpdateDatabaseStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UpdateDatabaseStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterUpdateDatabaseStmt(s)
	}
}

func (s *UpdateDatabaseStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitUpdateDatabaseStmt(s)
	}
}

func (s *UpdateDatabaseStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitUpdateDatabaseStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) UpdateDatabaseStmt() (localctx IUpdateDatabaseStmtContext) {
	localctx = NewUpdateDatabaseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, SQLParserRULE_updateDatabaseStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(252)
		p.Match(SQLParserT_UPDATE)
	}
	{
		p.SetState(253)
		p.Match(SQLParserT_DATASBAE)
	}
	{
		p.SetState(254)
		p.DatabaseName()
	}
	p.SetState(257)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_WITH {
		{
			p.SetState(255)
			p.Match(SQLParserT_WITH)
		}
		{
			p.SetState(256)
			p.WithClauseList()
		}

	}
	p.SetState(261)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_COMMA {
		{
			p.SetState(259)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(260)
			p.IntervalDefineList()
		}

	}

	return localctx
}

// IDropDatabaseStmtContext is an interface to support dynamic dispatch.
type IDropDatabaseStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDropDatabaseStmtContext differentiates from other interfaces.
	IsDropDatabaseStmtContext()
}

type DropDatabaseStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDropDatabaseStmtContext() *DropDatabaseStmtContext {
	var p = new(DropDatabaseStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_dropDatabaseStmt
	return p
}

func (*DropDatabaseStmtContext) IsDropDatabaseStmtContext() {}

func NewDropDatabaseStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DropDatabaseStmtContext {
	var p = new(DropDatabaseStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_dropDatabaseStmt

	return p
}

func (s *DropDatabaseStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *DropDatabaseStmtContext) T_DROP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DROP, 0)
}

func (s *DropDatabaseStmtContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *DropDatabaseStmtContext) DatabaseName() IDatabaseNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDatabaseNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDatabaseNameContext)
}

func (s *DropDatabaseStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DropDatabaseStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DropDatabaseStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDropDatabaseStmt(s)
	}
}

func (s *DropDatabaseStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDropDatabaseStmt(s)
	}
}

func (s *DropDatabaseStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDropDatabaseStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) DropDatabaseStmt() (localctx IDropDatabaseStmtContext) {
	localctx = NewDropDatabaseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, SQLParserRULE_dropDatabaseStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(263)
		p.Match(SQLParserT_DROP)
	}
	{
		p.SetState(264)
		p.Match(SQLParserT_DATASBAE)
	}
	{
		p.SetState(265)
		p.DatabaseName()
	}

	return localctx
}

// IShowDatabasesStmtContext is an interface to support dynamic dispatch.
type IShowDatabasesStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowDatabasesStmtContext differentiates from other interfaces.
	IsShowDatabasesStmtContext()
}

type ShowDatabasesStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowDatabasesStmtContext() *ShowDatabasesStmtContext {
	var p = new(ShowDatabasesStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showDatabasesStmt
	return p
}

func (*ShowDatabasesStmtContext) IsShowDatabasesStmtContext() {}

func NewShowDatabasesStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowDatabasesStmtContext {
	var p = new(ShowDatabasesStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showDatabasesStmt

	return p
}

func (s *ShowDatabasesStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowDatabasesStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowDatabasesStmtContext) T_DATASBAES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAES, 0)
}

func (s *ShowDatabasesStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowDatabasesStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowDatabasesStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowDatabasesStmt(s)
	}
}

func (s *ShowDatabasesStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowDatabasesStmt(s)
	}
}

func (s *ShowDatabasesStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowDatabasesStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowDatabasesStmt() (localctx IShowDatabasesStmtContext) {
	localctx = NewShowDatabasesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, SQLParserRULE_showDatabasesStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(267)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(268)
		p.Match(SQLParserT_DATASBAES)
	}

	return localctx
}

// IShowNodeStmtContext is an interface to support dynamic dispatch.
type IShowNodeStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowNodeStmtContext differentiates from other interfaces.
	IsShowNodeStmtContext()
}

type ShowNodeStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowNodeStmtContext() *ShowNodeStmtContext {
	var p = new(ShowNodeStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showNodeStmt
	return p
}

func (*ShowNodeStmtContext) IsShowNodeStmtContext() {}

func NewShowNodeStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowNodeStmtContext {
	var p = new(ShowNodeStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showNodeStmt

	return p
}

func (s *ShowNodeStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowNodeStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowNodeStmtContext) T_NODE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NODE, 0)
}

func (s *ShowNodeStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowNodeStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowNodeStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowNodeStmt(s)
	}
}

func (s *ShowNodeStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowNodeStmt(s)
	}
}

func (s *ShowNodeStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowNodeStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowNodeStmt() (localctx IShowNodeStmtContext) {
	localctx = NewShowNodeStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, SQLParserRULE_showNodeStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(270)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(271)
		p.Match(SQLParserT_NODE)
	}

	return localctx
}

// IShowMeasurementsStmtContext is an interface to support dynamic dispatch.
type IShowMeasurementsStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowMeasurementsStmtContext differentiates from other interfaces.
	IsShowMeasurementsStmtContext()
}

type ShowMeasurementsStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowMeasurementsStmtContext() *ShowMeasurementsStmtContext {
	var p = new(ShowMeasurementsStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showMeasurementsStmt
	return p
}

func (*ShowMeasurementsStmtContext) IsShowMeasurementsStmtContext() {}

func NewShowMeasurementsStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowMeasurementsStmtContext {
	var p = new(ShowMeasurementsStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showMeasurementsStmt

	return p
}

func (s *ShowMeasurementsStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowMeasurementsStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowMeasurementsStmtContext) T_MEASUREMENTS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MEASUREMENTS, 0)
}

func (s *ShowMeasurementsStmtContext) WithMeasurementClause() IWithMeasurementClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithMeasurementClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWithMeasurementClauseContext)
}

func (s *ShowMeasurementsStmtContext) LimitClause() ILimitClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimitClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ShowMeasurementsStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMeasurementsStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowMeasurementsStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowMeasurementsStmt(s)
	}
}

func (s *ShowMeasurementsStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowMeasurementsStmt(s)
	}
}

func (s *ShowMeasurementsStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowMeasurementsStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowMeasurementsStmt() (localctx IShowMeasurementsStmtContext) {
	localctx = NewShowMeasurementsStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, SQLParserRULE_showMeasurementsStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(273)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(274)
		p.Match(SQLParserT_MEASUREMENTS)
	}
	p.SetState(276)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_WITH {
		{
			p.SetState(275)
			p.WithMeasurementClause()
		}

	}
	p.SetState(279)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(278)
			p.LimitClause()
		}

	}

	return localctx
}

// IShowTagKeysStmtContext is an interface to support dynamic dispatch.
type IShowTagKeysStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowTagKeysStmtContext differentiates from other interfaces.
	IsShowTagKeysStmtContext()
}

type ShowTagKeysStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowTagKeysStmtContext() *ShowTagKeysStmtContext {
	var p = new(ShowTagKeysStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showTagKeysStmt
	return p
}

func (*ShowTagKeysStmtContext) IsShowTagKeysStmtContext() {}

func NewShowTagKeysStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowTagKeysStmtContext {
	var p = new(ShowTagKeysStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showTagKeysStmt

	return p
}

func (s *ShowTagKeysStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowTagKeysStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowTagKeysStmtContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *ShowTagKeysStmtContext) T_KEYS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEYS, 0)
}

func (s *ShowTagKeysStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowTagKeysStmtContext) MetricName() IMetricNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetricNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetricNameContext)
}

func (s *ShowTagKeysStmtContext) LimitClause() ILimitClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimitClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ShowTagKeysStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowTagKeysStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowTagKeysStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowTagKeysStmt(s)
	}
}

func (s *ShowTagKeysStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowTagKeysStmt(s)
	}
}

func (s *ShowTagKeysStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowTagKeysStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowTagKeysStmt() (localctx IShowTagKeysStmtContext) {
	localctx = NewShowTagKeysStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, SQLParserRULE_showTagKeysStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(281)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(282)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(283)
		p.Match(SQLParserT_KEYS)
	}
	{
		p.SetState(284)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(285)
		p.MetricName()
	}
	p.SetState(287)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(286)
			p.LimitClause()
		}

	}

	return localctx
}

// IShowInfoStmtContext is an interface to support dynamic dispatch.
type IShowInfoStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowInfoStmtContext differentiates from other interfaces.
	IsShowInfoStmtContext()
}

type ShowInfoStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowInfoStmtContext() *ShowInfoStmtContext {
	var p = new(ShowInfoStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showInfoStmt
	return p
}

func (*ShowInfoStmtContext) IsShowInfoStmtContext() {}

func NewShowInfoStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowInfoStmtContext {
	var p = new(ShowInfoStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showInfoStmt

	return p
}

func (s *ShowInfoStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowInfoStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowInfoStmtContext) T_INFO() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INFO, 0)
}

func (s *ShowInfoStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowInfoStmtContext) MetricName() IMetricNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetricNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetricNameContext)
}

func (s *ShowInfoStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowInfoStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowInfoStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowInfoStmt(s)
	}
}

func (s *ShowInfoStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowInfoStmt(s)
	}
}

func (s *ShowInfoStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowInfoStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowInfoStmt() (localctx IShowInfoStmtContext) {
	localctx = NewShowInfoStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, SQLParserRULE_showInfoStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(289)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(290)
		p.Match(SQLParserT_INFO)
	}
	{
		p.SetState(291)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(292)
		p.MetricName()
	}

	return localctx
}

// IShowTagValuesStmtContext is an interface to support dynamic dispatch.
type IShowTagValuesStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowTagValuesStmtContext differentiates from other interfaces.
	IsShowTagValuesStmtContext()
}

type ShowTagValuesStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowTagValuesStmtContext() *ShowTagValuesStmtContext {
	var p = new(ShowTagValuesStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showTagValuesStmt
	return p
}

func (*ShowTagValuesStmtContext) IsShowTagValuesStmtContext() {}

func NewShowTagValuesStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowTagValuesStmtContext {
	var p = new(ShowTagValuesStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showTagValuesStmt

	return p
}

func (s *ShowTagValuesStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowTagValuesStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowTagValuesStmtContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *ShowTagValuesStmtContext) T_VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUES, 0)
}

func (s *ShowTagValuesStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowTagValuesStmtContext) MetricName() IMetricNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetricNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetricNameContext)
}

func (s *ShowTagValuesStmtContext) WithTagClause() IWithTagClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithTagClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWithTagClauseContext)
}

func (s *ShowTagValuesStmtContext) WhereTagCascade() IWhereTagCascadeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWhereTagCascadeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWhereTagCascadeContext)
}

func (s *ShowTagValuesStmtContext) LimitClause() ILimitClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimitClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ShowTagValuesStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowTagValuesStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowTagValuesStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowTagValuesStmt(s)
	}
}

func (s *ShowTagValuesStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowTagValuesStmt(s)
	}
}

func (s *ShowTagValuesStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowTagValuesStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowTagValuesStmt() (localctx IShowTagValuesStmtContext) {
	localctx = NewShowTagValuesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, SQLParserRULE_showTagValuesStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(294)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(295)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(296)
		p.Match(SQLParserT_VALUES)
	}
	{
		p.SetState(297)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(298)
		p.MetricName()
	}
	{
		p.SetState(299)
		p.WithTagClause()
	}
	p.SetState(301)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_WHERE {
		{
			p.SetState(300)
			p.WhereTagCascade()
		}

	}
	p.SetState(304)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(303)
			p.LimitClause()
		}

	}

	return localctx
}

// IShowTagValuesInfoStmtContext is an interface to support dynamic dispatch.
type IShowTagValuesInfoStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowTagValuesInfoStmtContext differentiates from other interfaces.
	IsShowTagValuesInfoStmtContext()
}

type ShowTagValuesInfoStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowTagValuesInfoStmtContext() *ShowTagValuesInfoStmtContext {
	var p = new(ShowTagValuesInfoStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showTagValuesInfoStmt
	return p
}

func (*ShowTagValuesInfoStmtContext) IsShowTagValuesInfoStmtContext() {}

func NewShowTagValuesInfoStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowTagValuesInfoStmtContext {
	var p = new(ShowTagValuesInfoStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showTagValuesInfoStmt

	return p
}

func (s *ShowTagValuesInfoStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowTagValuesInfoStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowTagValuesInfoStmtContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *ShowTagValuesInfoStmtContext) T_VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUES, 0)
}

func (s *ShowTagValuesInfoStmtContext) T_INFO() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INFO, 0)
}

func (s *ShowTagValuesInfoStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowTagValuesInfoStmtContext) MetricName() IMetricNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetricNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetricNameContext)
}

func (s *ShowTagValuesInfoStmtContext) WithTagClause() IWithTagClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithTagClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWithTagClauseContext)
}

func (s *ShowTagValuesInfoStmtContext) WhereTagCascade() IWhereTagCascadeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWhereTagCascadeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWhereTagCascadeContext)
}

func (s *ShowTagValuesInfoStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowTagValuesInfoStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowTagValuesInfoStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowTagValuesInfoStmt(s)
	}
}

func (s *ShowTagValuesInfoStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowTagValuesInfoStmt(s)
	}
}

func (s *ShowTagValuesInfoStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowTagValuesInfoStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowTagValuesInfoStmt() (localctx IShowTagValuesInfoStmtContext) {
	localctx = NewShowTagValuesInfoStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, SQLParserRULE_showTagValuesInfoStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(306)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(307)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(308)
		p.Match(SQLParserT_VALUES)
	}
	{
		p.SetState(309)
		p.Match(SQLParserT_INFO)
	}
	{
		p.SetState(310)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(311)
		p.MetricName()
	}
	{
		p.SetState(312)
		p.WithTagClause()
	}
	{
		p.SetState(313)
		p.WhereTagCascade()
	}

	return localctx
}

// IShowFieldKeysStmtContext is an interface to support dynamic dispatch.
type IShowFieldKeysStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowFieldKeysStmtContext differentiates from other interfaces.
	IsShowFieldKeysStmtContext()
}

type ShowFieldKeysStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowFieldKeysStmtContext() *ShowFieldKeysStmtContext {
	var p = new(ShowFieldKeysStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showFieldKeysStmt
	return p
}

func (*ShowFieldKeysStmtContext) IsShowFieldKeysStmtContext() {}

func NewShowFieldKeysStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowFieldKeysStmtContext {
	var p = new(ShowFieldKeysStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showFieldKeysStmt

	return p
}

func (s *ShowFieldKeysStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowFieldKeysStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowFieldKeysStmtContext) T_FIELD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIELD, 0)
}

func (s *ShowFieldKeysStmtContext) T_KEYS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEYS, 0)
}

func (s *ShowFieldKeysStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowFieldKeysStmtContext) MetricName() IMetricNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetricNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetricNameContext)
}

func (s *ShowFieldKeysStmtContext) LimitClause() ILimitClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimitClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ShowFieldKeysStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowFieldKeysStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowFieldKeysStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowFieldKeysStmt(s)
	}
}

func (s *ShowFieldKeysStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowFieldKeysStmt(s)
	}
}

func (s *ShowFieldKeysStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowFieldKeysStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowFieldKeysStmt() (localctx IShowFieldKeysStmtContext) {
	localctx = NewShowFieldKeysStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, SQLParserRULE_showFieldKeysStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(315)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(316)
		p.Match(SQLParserT_FIELD)
	}
	{
		p.SetState(317)
		p.Match(SQLParserT_KEYS)
	}
	{
		p.SetState(318)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(319)
		p.MetricName()
	}
	p.SetState(321)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(320)
			p.LimitClause()
		}

	}

	return localctx
}

// IShowQueriesStmtContext is an interface to support dynamic dispatch.
type IShowQueriesStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowQueriesStmtContext differentiates from other interfaces.
	IsShowQueriesStmtContext()
}

type ShowQueriesStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowQueriesStmtContext() *ShowQueriesStmtContext {
	var p = new(ShowQueriesStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showQueriesStmt
	return p
}

func (*ShowQueriesStmtContext) IsShowQueriesStmtContext() {}

func NewShowQueriesStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowQueriesStmtContext {
	var p = new(ShowQueriesStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showQueriesStmt

	return p
}

func (s *ShowQueriesStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowQueriesStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowQueriesStmtContext) T_QUERIES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUERIES, 0)
}

func (s *ShowQueriesStmtContext) LimitClause() ILimitClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimitClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ShowQueriesStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowQueriesStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowQueriesStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowQueriesStmt(s)
	}
}

func (s *ShowQueriesStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowQueriesStmt(s)
	}
}

func (s *ShowQueriesStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowQueriesStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowQueriesStmt() (localctx IShowQueriesStmtContext) {
	localctx = NewShowQueriesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, SQLParserRULE_showQueriesStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(323)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(324)
		p.Match(SQLParserT_QUERIES)
	}
	p.SetState(326)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(325)
			p.LimitClause()
		}

	}

	return localctx
}

// IShowStatsStmtContext is an interface to support dynamic dispatch.
type IShowStatsStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowStatsStmtContext differentiates from other interfaces.
	IsShowStatsStmtContext()
}

type ShowStatsStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowStatsStmtContext() *ShowStatsStmtContext {
	var p = new(ShowStatsStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showStatsStmt
	return p
}

func (*ShowStatsStmtContext) IsShowStatsStmtContext() {}

func NewShowStatsStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowStatsStmtContext {
	var p = new(ShowStatsStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showStatsStmt

	return p
}

func (s *ShowStatsStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowStatsStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowStatsStmtContext) T_STATS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATS, 0)
}

func (s *ShowStatsStmtContext) T_FOR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FOR, 0)
}

func (s *ShowStatsStmtContext) Module() IModuleContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IModuleContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IModuleContext)
}

func (s *ShowStatsStmtContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *ShowStatsStmtContext) Component() IComponentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IComponentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IComponentContext)
}

func (s *ShowStatsStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStatsStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowStatsStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowStatsStmt(s)
	}
}

func (s *ShowStatsStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowStatsStmt(s)
	}
}

func (s *ShowStatsStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowStatsStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowStatsStmt() (localctx IShowStatsStmtContext) {
	localctx = NewShowStatsStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, SQLParserRULE_showStatsStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(328)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(329)
		p.Match(SQLParserT_STATS)
	}
	p.SetState(332)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_FOR {
		{
			p.SetState(330)
			p.Match(SQLParserT_FOR)
		}
		{
			p.SetState(331)
			p.Module()
		}

	}
	p.SetState(336)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_WITH {
		{
			p.SetState(334)
			p.Match(SQLParserT_WITH)
		}
		{
			p.SetState(335)
			p.Component()
		}

	}

	return localctx
}

// IWithMeasurementClauseContext is an interface to support dynamic dispatch.
type IWithMeasurementClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWithMeasurementClauseContext differentiates from other interfaces.
	IsWithMeasurementClauseContext()
}

type WithMeasurementClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithMeasurementClauseContext() *WithMeasurementClauseContext {
	var p = new(WithMeasurementClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_withMeasurementClause
	return p
}

func (*WithMeasurementClauseContext) IsWithMeasurementClauseContext() {}

func NewWithMeasurementClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithMeasurementClauseContext {
	var p = new(WithMeasurementClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_withMeasurementClause

	return p
}

func (s *WithMeasurementClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *WithMeasurementClauseContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *WithMeasurementClauseContext) T_MEASUREMENT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MEASUREMENT, 0)
}

func (s *WithMeasurementClauseContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *WithMeasurementClauseContext) MetricName() IMetricNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetricNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetricNameContext)
}

func (s *WithMeasurementClauseContext) T_REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REGEXP, 0)
}

func (s *WithMeasurementClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithMeasurementClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WithMeasurementClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWithMeasurementClause(s)
	}
}

func (s *WithMeasurementClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWithMeasurementClause(s)
	}
}

func (s *WithMeasurementClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWithMeasurementClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) WithMeasurementClause() (localctx IWithMeasurementClauseContext) {
	localctx = NewWithMeasurementClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, SQLParserRULE_withMeasurementClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(338)
		p.Match(SQLParserT_WITH)
	}
	{
		p.SetState(339)
		p.Match(SQLParserT_MEASUREMENT)
	}
	p.SetState(344)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_EQUAL:
		{
			p.SetState(340)
			p.Match(SQLParserT_EQUAL)
		}
		{
			p.SetState(341)
			p.MetricName()
		}

	case SQLParserT_REGEXP:
		{
			p.SetState(342)
			p.Match(SQLParserT_REGEXP)
		}
		{
			p.SetState(343)
			p.MetricName()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IWithTagClauseContext is an interface to support dynamic dispatch.
type IWithTagClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWithTagClauseContext differentiates from other interfaces.
	IsWithTagClauseContext()
}

type WithTagClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithTagClauseContext() *WithTagClauseContext {
	var p = new(WithTagClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_withTagClause
	return p
}

func (*WithTagClauseContext) IsWithTagClauseContext() {}

func NewWithTagClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithTagClauseContext {
	var p = new(WithTagClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_withTagClause

	return p
}

func (s *WithTagClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *WithTagClauseContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *WithTagClauseContext) T_KEY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEY, 0)
}

func (s *WithTagClauseContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *WithTagClauseContext) TagKey() ITagKeyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagKeyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagKeyContext)
}

func (s *WithTagClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithTagClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WithTagClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWithTagClause(s)
	}
}

func (s *WithTagClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWithTagClause(s)
	}
}

func (s *WithTagClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWithTagClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) WithTagClause() (localctx IWithTagClauseContext) {
	localctx = NewWithTagClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, SQLParserRULE_withTagClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(346)
		p.Match(SQLParserT_WITH)
	}
	{
		p.SetState(347)
		p.Match(SQLParserT_KEY)
	}
	{
		p.SetState(348)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(349)
		p.TagKey()
	}

	return localctx
}

// IWhereTagCascadeContext is an interface to support dynamic dispatch.
type IWhereTagCascadeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWhereTagCascadeContext differentiates from other interfaces.
	IsWhereTagCascadeContext()
}

type WhereTagCascadeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhereTagCascadeContext() *WhereTagCascadeContext {
	var p = new(WhereTagCascadeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_whereTagCascade
	return p
}

func (*WhereTagCascadeContext) IsWhereTagCascadeContext() {}

func NewWhereTagCascadeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereTagCascadeContext {
	var p = new(WhereTagCascadeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_whereTagCascade

	return p
}

func (s *WhereTagCascadeContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereTagCascadeContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *WhereTagCascadeContext) TagCascadeExpr() ITagCascadeExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagCascadeExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagCascadeExprContext)
}

func (s *WhereTagCascadeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereTagCascadeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WhereTagCascadeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWhereTagCascade(s)
	}
}

func (s *WhereTagCascadeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWhereTagCascade(s)
	}
}

func (s *WhereTagCascadeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWhereTagCascade(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) WhereTagCascade() (localctx IWhereTagCascadeContext) {
	localctx = NewWhereTagCascadeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, SQLParserRULE_whereTagCascade)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(351)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(352)
		p.TagCascadeExpr()
	}

	return localctx
}

// IKillQueryStmtContext is an interface to support dynamic dispatch.
type IKillQueryStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsKillQueryStmtContext differentiates from other interfaces.
	IsKillQueryStmtContext()
}

type KillQueryStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyKillQueryStmtContext() *KillQueryStmtContext {
	var p = new(KillQueryStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_killQueryStmt
	return p
}

func (*KillQueryStmtContext) IsKillQueryStmtContext() {}

func NewKillQueryStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *KillQueryStmtContext {
	var p = new(KillQueryStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_killQueryStmt

	return p
}

func (s *KillQueryStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *KillQueryStmtContext) T_KILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KILL, 0)
}

func (s *KillQueryStmtContext) T_QUERY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUERY, 0)
}

func (s *KillQueryStmtContext) QueryId() IQueryIdContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQueryIdContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IQueryIdContext)
}

func (s *KillQueryStmtContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *KillQueryStmtContext) ServerId() IServerIdContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IServerIdContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IServerIdContext)
}

func (s *KillQueryStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *KillQueryStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *KillQueryStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterKillQueryStmt(s)
	}
}

func (s *KillQueryStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitKillQueryStmt(s)
	}
}

func (s *KillQueryStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitKillQueryStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) KillQueryStmt() (localctx IKillQueryStmtContext) {
	localctx = NewKillQueryStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, SQLParserRULE_killQueryStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(354)
		p.Match(SQLParserT_KILL)
	}
	{
		p.SetState(355)
		p.Match(SQLParserT_QUERY)
	}
	{
		p.SetState(356)
		p.QueryId()
	}
	p.SetState(359)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_ON {
		{
			p.SetState(357)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(358)
			p.ServerId()
		}

	}

	return localctx
}

// IQueryIdContext is an interface to support dynamic dispatch.
type IQueryIdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQueryIdContext differentiates from other interfaces.
	IsQueryIdContext()
}

type QueryIdContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryIdContext() *QueryIdContext {
	var p = new(QueryIdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_queryId
	return p
}

func (*QueryIdContext) IsQueryIdContext() {}

func NewQueryIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryIdContext {
	var p = new(QueryIdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_queryId

	return p
}

func (s *QueryIdContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryIdContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *QueryIdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryIdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryIdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterQueryId(s)
	}
}

func (s *QueryIdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitQueryId(s)
	}
}

func (s *QueryIdContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitQueryId(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) QueryId() (localctx IQueryIdContext) {
	localctx = NewQueryIdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, SQLParserRULE_queryId)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(361)
		p.Match(SQLParserL_INT)
	}

	return localctx
}

// IServerIdContext is an interface to support dynamic dispatch.
type IServerIdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsServerIdContext differentiates from other interfaces.
	IsServerIdContext()
}

type ServerIdContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyServerIdContext() *ServerIdContext {
	var p = new(ServerIdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_serverId
	return p
}

func (*ServerIdContext) IsServerIdContext() {}

func NewServerIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ServerIdContext {
	var p = new(ServerIdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_serverId

	return p
}

func (s *ServerIdContext) GetParser() antlr.Parser { return s.parser }

func (s *ServerIdContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *ServerIdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ServerIdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ServerIdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterServerId(s)
	}
}

func (s *ServerIdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitServerId(s)
	}
}

func (s *ServerIdContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitServerId(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ServerId() (localctx IServerIdContext) {
	localctx = NewServerIdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, SQLParserRULE_serverId)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(363)
		p.Match(SQLParserL_INT)
	}

	return localctx
}

// IModuleContext is an interface to support dynamic dispatch.
type IModuleContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsModuleContext differentiates from other interfaces.
	IsModuleContext()
}

type ModuleContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModuleContext() *ModuleContext {
	var p = new(ModuleContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_module
	return p
}

func (*ModuleContext) IsModuleContext() {}

func NewModuleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModuleContext {
	var p = new(ModuleContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_module

	return p
}

func (s *ModuleContext) GetParser() antlr.Parser { return s.parser }

func (s *ModuleContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ModuleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ModuleContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ModuleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterModule(s)
	}
}

func (s *ModuleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitModule(s)
	}
}

func (s *ModuleContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitModule(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Module() (localctx IModuleContext) {
	localctx = NewModuleContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, SQLParserRULE_module)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(365)
		p.Ident()
	}

	return localctx
}

// IComponentContext is an interface to support dynamic dispatch.
type IComponentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsComponentContext differentiates from other interfaces.
	IsComponentContext()
}

type ComponentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComponentContext() *ComponentContext {
	var p = new(ComponentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_component
	return p
}

func (*ComponentContext) IsComponentContext() {}

func NewComponentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComponentContext {
	var p = new(ComponentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_component

	return p
}

func (s *ComponentContext) GetParser() antlr.Parser { return s.parser }

func (s *ComponentContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ComponentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComponentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComponentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterComponent(s)
	}
}

func (s *ComponentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitComponent(s)
	}
}

func (s *ComponentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitComponent(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Component() (localctx IComponentContext) {
	localctx = NewComponentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, SQLParserRULE_component)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(367)
		p.Ident()
	}

	return localctx
}

// IQueryStmtContext is an interface to support dynamic dispatch.
type IQueryStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQueryStmtContext differentiates from other interfaces.
	IsQueryStmtContext()
}

type QueryStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryStmtContext() *QueryStmtContext {
	var p = new(QueryStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_queryStmt
	return p
}

func (*QueryStmtContext) IsQueryStmtContext() {}

func NewQueryStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryStmtContext {
	var p = new(QueryStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_queryStmt

	return p
}

func (s *QueryStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryStmtContext) T_SELECT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SELECT, 0)
}

func (s *QueryStmtContext) Fields() IFieldsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldsContext)
}

func (s *QueryStmtContext) FromClause() IFromClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFromClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFromClauseContext)
}

func (s *QueryStmtContext) T_EXPLAIN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EXPLAIN, 0)
}

func (s *QueryStmtContext) WhereClause() IWhereClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWhereClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWhereClauseContext)
}

func (s *QueryStmtContext) GroupByClause() IGroupByClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IGroupByClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IGroupByClauseContext)
}

func (s *QueryStmtContext) IntervalByClause() IIntervalByClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntervalByClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntervalByClauseContext)
}

func (s *QueryStmtContext) OrderByClause() IOrderByClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOrderByClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOrderByClauseContext)
}

func (s *QueryStmtContext) LimitClause() ILimitClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimitClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *QueryStmtContext) T_WITH_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH_VALUE, 0)
}

func (s *QueryStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterQueryStmt(s)
	}
}

func (s *QueryStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitQueryStmt(s)
	}
}

func (s *QueryStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitQueryStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) QueryStmt() (localctx IQueryStmtContext) {
	localctx = NewQueryStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, SQLParserRULE_queryStmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(370)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_EXPLAIN {
		{
			p.SetState(369)
			p.Match(SQLParserT_EXPLAIN)
		}

	}
	{
		p.SetState(372)
		p.Match(SQLParserT_SELECT)
	}
	{
		p.SetState(373)
		p.Fields()
	}
	{
		p.SetState(374)
		p.FromClause()
	}
	p.SetState(376)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_WHERE {
		{
			p.SetState(375)
			p.WhereClause()
		}

	}
	p.SetState(379)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_GROUP {
		{
			p.SetState(378)
			p.GroupByClause()
		}

	}
	p.SetState(382)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_INTERVAL {
		{
			p.SetState(381)
			p.IntervalByClause()
		}

	}
	p.SetState(385)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_ORDER {
		{
			p.SetState(384)
			p.OrderByClause()
		}

	}
	p.SetState(388)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(387)
			p.LimitClause()
		}

	}
	p.SetState(391)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_WITH_VALUE {
		{
			p.SetState(390)
			p.Match(SQLParserT_WITH_VALUE)
		}

	}

	return localctx
}

// IFieldsContext is an interface to support dynamic dispatch.
type IFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldsContext differentiates from other interfaces.
	IsFieldsContext()
}

type FieldsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldsContext() *FieldsContext {
	var p = new(FieldsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_fields
	return p
}

func (*FieldsContext) IsFieldsContext() {}

func NewFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldsContext {
	var p = new(FieldsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_fields

	return p
}

func (s *FieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldsContext) AllField() []IFieldContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFieldContext)(nil)).Elem())
	var tst = make([]IFieldContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFieldContext)
		}
	}

	return tst
}

func (s *FieldsContext) Field(i int) IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *FieldsContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *FieldsContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *FieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFields(s)
	}
}

func (s *FieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFields(s)
	}
}

func (s *FieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFields(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Fields() (localctx IFieldsContext) {
	localctx = NewFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, SQLParserRULE_fields)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(393)
		p.Field()
	}
	p.SetState(398)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserT_COMMA {
		{
			p.SetState(394)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(395)
			p.Field()
		}

		p.SetState(400)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IFieldContext is an interface to support dynamic dispatch.
type IFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldContext differentiates from other interfaces.
	IsFieldContext()
}

type FieldContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldContext() *FieldContext {
	var p = new(FieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_field
	return p
}

func (*FieldContext) IsFieldContext() {}

func NewFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldContext {
	var p = new(FieldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_field

	return p
}

func (s *FieldContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *FieldContext) Alias() IAliasContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAliasContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAliasContext)
}

func (s *FieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterField(s)
	}
}

func (s *FieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitField(s)
	}
}

func (s *FieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitField(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Field() (localctx IFieldContext) {
	localctx = NewFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, SQLParserRULE_field)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(401)
		p.expr(0)
	}
	p.SetState(403)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_AS {
		{
			p.SetState(402)
			p.Alias()
		}

	}

	return localctx
}

// IAliasContext is an interface to support dynamic dispatch.
type IAliasContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAliasContext differentiates from other interfaces.
	IsAliasContext()
}

type AliasContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAliasContext() *AliasContext {
	var p = new(AliasContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_alias
	return p
}

func (*AliasContext) IsAliasContext() {}

func NewAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AliasContext {
	var p = new(AliasContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_alias

	return p
}

func (s *AliasContext) GetParser() antlr.Parser { return s.parser }

func (s *AliasContext) T_AS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AS, 0)
}

func (s *AliasContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *AliasContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AliasContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AliasContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterAlias(s)
	}
}

func (s *AliasContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitAlias(s)
	}
}

func (s *AliasContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitAlias(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Alias() (localctx IAliasContext) {
	localctx = NewAliasContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, SQLParserRULE_alias)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(405)
		p.Match(SQLParserT_AS)
	}
	{
		p.SetState(406)
		p.Ident()
	}

	return localctx
}

// IFromClauseContext is an interface to support dynamic dispatch.
type IFromClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFromClauseContext differentiates from other interfaces.
	IsFromClauseContext()
}

type FromClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFromClauseContext() *FromClauseContext {
	var p = new(FromClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_fromClause
	return p
}

func (*FromClauseContext) IsFromClauseContext() {}

func NewFromClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FromClauseContext {
	var p = new(FromClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_fromClause

	return p
}

func (s *FromClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *FromClauseContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *FromClauseContext) MetricName() IMetricNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetricNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetricNameContext)
}

func (s *FromClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FromClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FromClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFromClause(s)
	}
}

func (s *FromClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFromClause(s)
	}
}

func (s *FromClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFromClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) FromClause() (localctx IFromClauseContext) {
	localctx = NewFromClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, SQLParserRULE_fromClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(408)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(409)
		p.MetricName()
	}

	return localctx
}

// IWhereClauseContext is an interface to support dynamic dispatch.
type IWhereClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWhereClauseContext differentiates from other interfaces.
	IsWhereClauseContext()
}

type WhereClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhereClauseContext() *WhereClauseContext {
	var p = new(WhereClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_whereClause
	return p
}

func (*WhereClauseContext) IsWhereClauseContext() {}

func NewWhereClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereClauseContext {
	var p = new(WhereClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_whereClause

	return p
}

func (s *WhereClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereClauseContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *WhereClauseContext) ClauseBooleanExpr() IClauseBooleanExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IClauseBooleanExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IClauseBooleanExprContext)
}

func (s *WhereClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WhereClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWhereClause(s)
	}
}

func (s *WhereClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWhereClause(s)
	}
}

func (s *WhereClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWhereClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) WhereClause() (localctx IWhereClauseContext) {
	localctx = NewWhereClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, SQLParserRULE_whereClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(411)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(412)
		p.clauseBooleanExpr(0)
	}

	return localctx
}

// IClauseBooleanExprContext is an interface to support dynamic dispatch.
type IClauseBooleanExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsClauseBooleanExprContext differentiates from other interfaces.
	IsClauseBooleanExprContext()
}

type ClauseBooleanExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyClauseBooleanExprContext() *ClauseBooleanExprContext {
	var p = new(ClauseBooleanExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_clauseBooleanExpr
	return p
}

func (*ClauseBooleanExprContext) IsClauseBooleanExprContext() {}

func NewClauseBooleanExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClauseBooleanExprContext {
	var p = new(ClauseBooleanExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_clauseBooleanExpr

	return p
}

func (s *ClauseBooleanExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ClauseBooleanExprContext) TagBooleanExpr() ITagBooleanExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagBooleanExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagBooleanExprContext)
}

func (s *ClauseBooleanExprContext) TimeExpr() ITimeExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITimeExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITimeExprContext)
}

func (s *ClauseBooleanExprContext) AllClauseBooleanExpr() []IClauseBooleanExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IClauseBooleanExprContext)(nil)).Elem())
	var tst = make([]IClauseBooleanExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IClauseBooleanExprContext)
		}
	}

	return tst
}

func (s *ClauseBooleanExprContext) ClauseBooleanExpr(i int) IClauseBooleanExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IClauseBooleanExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IClauseBooleanExprContext)
}

func (s *ClauseBooleanExprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *ClauseBooleanExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClauseBooleanExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ClauseBooleanExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterClauseBooleanExpr(s)
	}
}

func (s *ClauseBooleanExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitClauseBooleanExpr(s)
	}
}

func (s *ClauseBooleanExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitClauseBooleanExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ClauseBooleanExpr() (localctx IClauseBooleanExprContext) {
	return p.clauseBooleanExpr(0)
}

func (p *SQLParser) clauseBooleanExpr(_p int) (localctx IClauseBooleanExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewClauseBooleanExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IClauseBooleanExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 82
	p.EnterRecursionRule(localctx, 82, SQLParserRULE_clauseBooleanExpr, _p)

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(417)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(415)
			p.tagBooleanExpr(0)
		}

	case 2:
		{
			p.SetState(416)
			p.TimeExpr()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(424)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewClauseBooleanExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_clauseBooleanExpr)
			p.SetState(419)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(420)
				p.Match(SQLParserT_AND)
			}
			{
				p.SetState(421)
				p.clauseBooleanExpr(2)
			}

		}
		p.SetState(426)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())
	}

	return localctx
}

// ITagCascadeExprContext is an interface to support dynamic dispatch.
type ITagCascadeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagCascadeExprContext differentiates from other interfaces.
	IsTagCascadeExprContext()
}

type TagCascadeExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagCascadeExprContext() *TagCascadeExprContext {
	var p = new(TagCascadeExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagCascadeExpr
	return p
}

func (*TagCascadeExprContext) IsTagCascadeExprContext() {}

func NewTagCascadeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagCascadeExprContext {
	var p = new(TagCascadeExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagCascadeExpr

	return p
}

func (s *TagCascadeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TagCascadeExprContext) TagEqualExpr() ITagEqualExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagEqualExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagEqualExprContext)
}

func (s *TagCascadeExprContext) TagBooleanExpr() ITagBooleanExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagBooleanExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagBooleanExprContext)
}

func (s *TagCascadeExprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *TagCascadeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagCascadeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TagCascadeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagCascadeExpr(s)
	}
}

func (s *TagCascadeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagCascadeExpr(s)
	}
}

func (s *TagCascadeExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagCascadeExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TagCascadeExpr() (localctx ITagCascadeExprContext) {
	localctx = NewTagCascadeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, SQLParserRULE_tagCascadeExpr)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(434)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 31, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(427)
			p.TagEqualExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(428)
			p.tagBooleanExpr(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(429)
			p.TagEqualExpr()
		}
		p.SetState(432)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserT_AND {
			{
				p.SetState(430)
				p.Match(SQLParserT_AND)
			}
			{
				p.SetState(431)
				p.tagBooleanExpr(0)
			}

		}

	}

	return localctx
}

// ITagEqualExprContext is an interface to support dynamic dispatch.
type ITagEqualExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagEqualExprContext differentiates from other interfaces.
	IsTagEqualExprContext()
}

type TagEqualExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagEqualExprContext() *TagEqualExprContext {
	var p = new(TagEqualExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagEqualExpr
	return p
}

func (*TagEqualExprContext) IsTagEqualExprContext() {}

func NewTagEqualExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagEqualExprContext {
	var p = new(TagEqualExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagEqualExpr

	return p
}

func (s *TagEqualExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TagEqualExprContext) T_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUE, 0)
}

func (s *TagEqualExprContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *TagEqualExprContext) TagValuePattern() ITagValuePatternContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagValuePatternContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagValuePatternContext)
}

func (s *TagEqualExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagEqualExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TagEqualExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagEqualExpr(s)
	}
}

func (s *TagEqualExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagEqualExpr(s)
	}
}

func (s *TagEqualExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagEqualExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TagEqualExpr() (localctx ITagEqualExprContext) {
	localctx = NewTagEqualExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, SQLParserRULE_tagEqualExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(436)
		p.Match(SQLParserT_VALUE)
	}
	{
		p.SetState(437)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(438)
		p.TagValuePattern()
	}

	return localctx
}

// ITagBooleanExprContext is an interface to support dynamic dispatch.
type ITagBooleanExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagBooleanExprContext differentiates from other interfaces.
	IsTagBooleanExprContext()
}

type TagBooleanExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagBooleanExprContext() *TagBooleanExprContext {
	var p = new(TagBooleanExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagBooleanExpr
	return p
}

func (*TagBooleanExprContext) IsTagBooleanExprContext() {}

func NewTagBooleanExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagBooleanExprContext {
	var p = new(TagBooleanExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagBooleanExpr

	return p
}

func (s *TagBooleanExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TagBooleanExprContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *TagBooleanExprContext) AllTagBooleanExpr() []ITagBooleanExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITagBooleanExprContext)(nil)).Elem())
	var tst = make([]ITagBooleanExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITagBooleanExprContext)
		}
	}

	return tst
}

func (s *TagBooleanExprContext) TagBooleanExpr(i int) ITagBooleanExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagBooleanExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITagBooleanExprContext)
}

func (s *TagBooleanExprContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *TagBooleanExprContext) TagKey() ITagKeyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagKeyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagKeyContext)
}

func (s *TagBooleanExprContext) TagValue() ITagValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagValueContext)
}

func (s *TagBooleanExprContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *TagBooleanExprContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *TagBooleanExprContext) T_REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REGEXP, 0)
}

func (s *TagBooleanExprContext) T_NOTEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL, 0)
}

func (s *TagBooleanExprContext) T_NOTEQUAL2() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL2, 0)
}

func (s *TagBooleanExprContext) TagValueList() ITagValueListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagValueListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagValueListContext)
}

func (s *TagBooleanExprContext) T_IN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IN, 0)
}

func (s *TagBooleanExprContext) T_NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOT, 0)
}

func (s *TagBooleanExprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *TagBooleanExprContext) T_OR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OR, 0)
}

func (s *TagBooleanExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagBooleanExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TagBooleanExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagBooleanExpr(s)
	}
}

func (s *TagBooleanExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagBooleanExpr(s)
	}
}

func (s *TagBooleanExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagBooleanExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TagBooleanExpr() (localctx ITagBooleanExprContext) {
	return p.tagBooleanExpr(0)
}

func (p *SQLParser) tagBooleanExpr(_p int) (localctx ITagBooleanExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewTagBooleanExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ITagBooleanExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 88
	p.EnterRecursionRule(localctx, 88, SQLParserRULE_tagBooleanExpr, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(459)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(441)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(442)
			p.tagBooleanExpr(0)
		}
		{
			p.SetState(443)
			p.Match(SQLParserT_CLOSE_P)
		}

	case 2:
		{
			p.SetState(445)
			p.TagKey()
		}
		{
			p.SetState(446)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-46)&-(0x1f+1)) == 0 && ((1<<uint((_la-46)))&((1<<(SQLParserT_LIKE-46))|(1<<(SQLParserT_EQUAL-46))|(1<<(SQLParserT_NOTEQUAL-46))|(1<<(SQLParserT_NOTEQUAL2-46))|(1<<(SQLParserT_REGEXP-46)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(447)
			p.TagValue()
		}

	case 3:
		{
			p.SetState(449)
			p.TagKey()
		}
		p.SetState(453)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_IN:
			{
				p.SetState(450)
				p.Match(SQLParserT_IN)
			}

		case SQLParserT_NOT:
			{
				p.SetState(451)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(452)
				p.Match(SQLParserT_IN)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(455)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(456)
			p.TagValueList()
		}
		{
			p.SetState(457)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(466)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewTagBooleanExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_tagBooleanExpr)
			p.SetState(461)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(462)
				_la = p.GetTokenStream().LA(1)

				if !(_la == SQLParserT_AND || _la == SQLParserT_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(463)
				p.tagBooleanExpr(2)
			}

		}
		p.SetState(468)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())
	}

	return localctx
}

// ITagValueListContext is an interface to support dynamic dispatch.
type ITagValueListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagValueListContext differentiates from other interfaces.
	IsTagValueListContext()
}

type TagValueListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagValueListContext() *TagValueListContext {
	var p = new(TagValueListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagValueList
	return p
}

func (*TagValueListContext) IsTagValueListContext() {}

func NewTagValueListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagValueListContext {
	var p = new(TagValueListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagValueList

	return p
}

func (s *TagValueListContext) GetParser() antlr.Parser { return s.parser }

func (s *TagValueListContext) AllTagValue() []ITagValueContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITagValueContext)(nil)).Elem())
	var tst = make([]ITagValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITagValueContext)
		}
	}

	return tst
}

func (s *TagValueListContext) TagValue(i int) ITagValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagValueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITagValueContext)
}

func (s *TagValueListContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *TagValueListContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *TagValueListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagValueListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TagValueListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagValueList(s)
	}
}

func (s *TagValueListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagValueList(s)
	}
}

func (s *TagValueListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagValueList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TagValueList() (localctx ITagValueListContext) {
	localctx = NewTagValueListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, SQLParserRULE_tagValueList)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(469)
		p.TagValue()
	}
	p.SetState(474)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserT_COMMA {
		{
			p.SetState(470)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(471)
			p.TagValue()
		}

		p.SetState(476)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ITimeExprContext is an interface to support dynamic dispatch.
type ITimeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTimeExprContext differentiates from other interfaces.
	IsTimeExprContext()
}

type TimeExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeExprContext() *TimeExprContext {
	var p = new(TimeExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_timeExpr
	return p
}

func (*TimeExprContext) IsTimeExprContext() {}

func NewTimeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeExprContext {
	var p = new(TimeExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_timeExpr

	return p
}

func (s *TimeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeExprContext) AllTimeBooleanExpr() []ITimeBooleanExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITimeBooleanExprContext)(nil)).Elem())
	var tst = make([]ITimeBooleanExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITimeBooleanExprContext)
		}
	}

	return tst
}

func (s *TimeExprContext) TimeBooleanExpr(i int) ITimeBooleanExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITimeBooleanExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITimeBooleanExprContext)
}

func (s *TimeExprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *TimeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTimeExpr(s)
	}
}

func (s *TimeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTimeExpr(s)
	}
}

func (s *TimeExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTimeExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TimeExpr() (localctx ITimeExprContext) {
	localctx = NewTimeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, SQLParserRULE_timeExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(477)
		p.TimeBooleanExpr()
	}
	p.SetState(480)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(478)
			p.Match(SQLParserT_AND)
		}
		{
			p.SetState(479)
			p.TimeBooleanExpr()
		}

	}

	return localctx
}

// ITimeBooleanExprContext is an interface to support dynamic dispatch.
type ITimeBooleanExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTimeBooleanExprContext differentiates from other interfaces.
	IsTimeBooleanExprContext()
}

type TimeBooleanExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeBooleanExprContext() *TimeBooleanExprContext {
	var p = new(TimeBooleanExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_timeBooleanExpr
	return p
}

func (*TimeBooleanExprContext) IsTimeBooleanExprContext() {}

func NewTimeBooleanExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeBooleanExprContext {
	var p = new(TimeBooleanExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_timeBooleanExpr

	return p
}

func (s *TimeBooleanExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeBooleanExprContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *TimeBooleanExprContext) BoolExprBinaryOperator() IBoolExprBinaryOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBoolExprBinaryOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBoolExprBinaryOperatorContext)
}

func (s *TimeBooleanExprContext) NowExpr() INowExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INowExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INowExprContext)
}

func (s *TimeBooleanExprContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *TimeBooleanExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeBooleanExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeBooleanExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTimeBooleanExpr(s)
	}
}

func (s *TimeBooleanExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTimeBooleanExpr(s)
	}
}

func (s *TimeBooleanExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTimeBooleanExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TimeBooleanExpr() (localctx ITimeBooleanExprContext) {
	localctx = NewTimeBooleanExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, SQLParserRULE_timeBooleanExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(482)
		p.Match(SQLParserT_TIME)
	}
	{
		p.SetState(483)
		p.BoolExprBinaryOperator()
	}
	p.SetState(486)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_NOW:
		{
			p.SetState(484)
			p.NowExpr()
		}

	case SQLParserT_CREATE, SQLParserT_INTERVAL, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NODE, SQLParserT_MEASUREMENTS, SQLParserT_MEASUREMENT, SQLParserT_FIELD, SQLParserT_TAG, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_PROFILE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR, SQLParserL_ID:
		{
			p.SetState(485)
			p.Ident()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// INowExprContext is an interface to support dynamic dispatch.
type INowExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNowExprContext differentiates from other interfaces.
	IsNowExprContext()
}

type NowExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNowExprContext() *NowExprContext {
	var p = new(NowExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_nowExpr
	return p
}

func (*NowExprContext) IsNowExprContext() {}

func NewNowExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NowExprContext {
	var p = new(NowExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_nowExpr

	return p
}

func (s *NowExprContext) GetParser() antlr.Parser { return s.parser }

func (s *NowExprContext) NowFunc() INowFuncContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INowFuncContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INowFuncContext)
}

func (s *NowExprContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *NowExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NowExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NowExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNowExpr(s)
	}
}

func (s *NowExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNowExpr(s)
	}
}

func (s *NowExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitNowExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) NowExpr() (localctx INowExprContext) {
	localctx = NewNowExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, SQLParserRULE_nowExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(488)
		p.NowFunc()
	}
	p.SetState(490)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(489)
			p.DurationLit()
		}

	}

	return localctx
}

// INowFuncContext is an interface to support dynamic dispatch.
type INowFuncContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNowFuncContext differentiates from other interfaces.
	IsNowFuncContext()
}

type NowFuncContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNowFuncContext() *NowFuncContext {
	var p = new(NowFuncContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_nowFunc
	return p
}

func (*NowFuncContext) IsNowFuncContext() {}

func NewNowFuncContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NowFuncContext {
	var p = new(NowFuncContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_nowFunc

	return p
}

func (s *NowFuncContext) GetParser() antlr.Parser { return s.parser }

func (s *NowFuncContext) T_NOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOW, 0)
}

func (s *NowFuncContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *NowFuncContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *NowFuncContext) ExprFuncParams() IExprFuncParamsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprFuncParamsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprFuncParamsContext)
}

func (s *NowFuncContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NowFuncContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NowFuncContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNowFunc(s)
	}
}

func (s *NowFuncContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNowFunc(s)
	}
}

func (s *NowFuncContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitNowFunc(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) NowFunc() (localctx INowFuncContext) {
	localctx = NewNowFuncContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, SQLParserRULE_nowFunc)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(492)
		p.Match(SQLParserT_NOW)
	}
	{
		p.SetState(493)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(495)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SQLParserT_CREATE)|(1<<SQLParserT_INTERVAL)|(1<<SQLParserT_SHARD)|(1<<SQLParserT_REPLICATION)|(1<<SQLParserT_TTL)|(1<<SQLParserT_KILL)|(1<<SQLParserT_ON)|(1<<SQLParserT_SHOW)|(1<<SQLParserT_DATASBAE)|(1<<SQLParserT_DATASBAES)|(1<<SQLParserT_NODE)|(1<<SQLParserT_MEASUREMENTS)|(1<<SQLParserT_MEASUREMENT)|(1<<SQLParserT_FIELD)|(1<<SQLParserT_TAG)|(1<<SQLParserT_KEYS)|(1<<SQLParserT_KEY)|(1<<SQLParserT_WITH)|(1<<SQLParserT_VALUES)|(1<<SQLParserT_FROM)|(1<<SQLParserT_WHERE)|(1<<SQLParserT_LIMIT))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(SQLParserT_QUERIES-32))|(1<<(SQLParserT_QUERY-32))|(1<<(SQLParserT_SELECT-32))|(1<<(SQLParserT_AS-32))|(1<<(SQLParserT_AND-32))|(1<<(SQLParserT_OR-32))|(1<<(SQLParserT_FILL-32))|(1<<(SQLParserT_NULL-32))|(1<<(SQLParserT_PREVIOUS-32))|(1<<(SQLParserT_ORDER-32))|(1<<(SQLParserT_ASC-32))|(1<<(SQLParserT_DESC-32))|(1<<(SQLParserT_LIKE-32))|(1<<(SQLParserT_NOT-32))|(1<<(SQLParserT_BETWEEN-32))|(1<<(SQLParserT_IS-32))|(1<<(SQLParserT_GROUP-32))|(1<<(SQLParserT_BY-32))|(1<<(SQLParserT_FOR-32))|(1<<(SQLParserT_STATS-32))|(1<<(SQLParserT_TIME-32))|(1<<(SQLParserT_PROFILE-32))|(1<<(SQLParserT_SECOND-32))|(1<<(SQLParserT_MINUTE-32))|(1<<(SQLParserT_HOUR-32))|(1<<(SQLParserT_DAY-32)))) != 0) || (((_la-64)&-(0x1f+1)) == 0 && ((1<<uint((_la-64)))&((1<<(SQLParserT_WEEK-64))|(1<<(SQLParserT_MONTH-64))|(1<<(SQLParserT_YEAR-64))|(1<<(SQLParserT_OPEN_P-64))|(1<<(SQLParserT_ADD-64))|(1<<(SQLParserT_SUB-64))|(1<<(SQLParserL_ID-64))|(1<<(SQLParserL_INT-64))|(1<<(SQLParserL_DEC-64)))) != 0) {
		{
			p.SetState(494)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(497)
		p.Match(SQLParserT_CLOSE_P)
	}

	return localctx
}

// IGroupByClauseContext is an interface to support dynamic dispatch.
type IGroupByClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGroupByClauseContext differentiates from other interfaces.
	IsGroupByClauseContext()
}

type GroupByClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByClauseContext() *GroupByClauseContext {
	var p = new(GroupByClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_groupByClause
	return p
}

func (*GroupByClauseContext) IsGroupByClauseContext() {}

func NewGroupByClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByClauseContext {
	var p = new(GroupByClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_groupByClause

	return p
}

func (s *GroupByClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByClauseContext) T_GROUP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GROUP, 0)
}

func (s *GroupByClauseContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *GroupByClauseContext) Dimensions() IDimensionsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDimensionsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDimensionsContext)
}

func (s *GroupByClauseContext) T_FILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FILL, 0)
}

func (s *GroupByClauseContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *GroupByClauseContext) FillOption() IFillOptionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFillOptionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFillOptionContext)
}

func (s *GroupByClauseContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *GroupByClauseContext) HavingClause() IHavingClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHavingClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHavingClauseContext)
}

func (s *GroupByClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupByClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterGroupByClause(s)
	}
}

func (s *GroupByClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitGroupByClause(s)
	}
}

func (s *GroupByClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitGroupByClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) GroupByClause() (localctx IGroupByClauseContext) {
	localctx = NewGroupByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, SQLParserRULE_groupByClause)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(499)
		p.Match(SQLParserT_GROUP)
	}
	{
		p.SetState(500)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(501)
		p.Dimensions()
	}
	p.SetState(507)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_FILL {
		{
			p.SetState(502)
			p.Match(SQLParserT_FILL)
		}
		{
			p.SetState(503)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(504)
			p.FillOption()
		}
		{
			p.SetState(505)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.SetState(510)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_HAVING {
		{
			p.SetState(509)
			p.HavingClause()
		}

	}

	return localctx
}

// IDimensionsContext is an interface to support dynamic dispatch.
type IDimensionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDimensionsContext differentiates from other interfaces.
	IsDimensionsContext()
}

type DimensionsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDimensionsContext() *DimensionsContext {
	var p = new(DimensionsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_dimensions
	return p
}

func (*DimensionsContext) IsDimensionsContext() {}

func NewDimensionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DimensionsContext {
	var p = new(DimensionsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_dimensions

	return p
}

func (s *DimensionsContext) GetParser() antlr.Parser { return s.parser }

func (s *DimensionsContext) AllDimension() []IDimensionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IDimensionContext)(nil)).Elem())
	var tst = make([]IDimensionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IDimensionContext)
		}
	}

	return tst
}

func (s *DimensionsContext) Dimension(i int) IDimensionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDimensionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IDimensionContext)
}

func (s *DimensionsContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *DimensionsContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *DimensionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DimensionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DimensionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDimensions(s)
	}
}

func (s *DimensionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDimensions(s)
	}
}

func (s *DimensionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDimensions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Dimensions() (localctx IDimensionsContext) {
	localctx = NewDimensionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, SQLParserRULE_dimensions)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(512)
		p.Dimension()
	}
	p.SetState(517)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserT_COMMA {
		{
			p.SetState(513)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(514)
			p.Dimension()
		}

		p.SetState(519)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IDimensionContext is an interface to support dynamic dispatch.
type IDimensionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDimensionContext differentiates from other interfaces.
	IsDimensionContext()
}

type DimensionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDimensionContext() *DimensionContext {
	var p = new(DimensionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_dimension
	return p
}

func (*DimensionContext) IsDimensionContext() {}

func NewDimensionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DimensionContext {
	var p = new(DimensionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_dimension

	return p
}

func (s *DimensionContext) GetParser() antlr.Parser { return s.parser }

func (s *DimensionContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *DimensionContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *DimensionContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *DimensionContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *DimensionContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *DimensionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DimensionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DimensionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDimension(s)
	}
}

func (s *DimensionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDimension(s)
	}
}

func (s *DimensionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDimension(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Dimension() (localctx IDimensionContext) {
	localctx = NewDimensionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, SQLParserRULE_dimension)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(526)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(520)
			p.Ident()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(521)
			p.Match(SQLParserT_TIME)
		}
		{
			p.SetState(522)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(523)
			p.DurationLit()
		}
		{
			p.SetState(524)
			p.Match(SQLParserT_CLOSE_P)
		}

	}

	return localctx
}

// IFillOptionContext is an interface to support dynamic dispatch.
type IFillOptionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFillOptionContext differentiates from other interfaces.
	IsFillOptionContext()
}

type FillOptionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFillOptionContext() *FillOptionContext {
	var p = new(FillOptionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_fillOption
	return p
}

func (*FillOptionContext) IsFillOptionContext() {}

func NewFillOptionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FillOptionContext {
	var p = new(FillOptionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_fillOption

	return p
}

func (s *FillOptionContext) GetParser() antlr.Parser { return s.parser }

func (s *FillOptionContext) T_NULL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NULL, 0)
}

func (s *FillOptionContext) T_PREVIOUS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PREVIOUS, 0)
}

func (s *FillOptionContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *FillOptionContext) L_DEC() antlr.TerminalNode {
	return s.GetToken(SQLParserL_DEC, 0)
}

func (s *FillOptionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillOptionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FillOptionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFillOption(s)
	}
}

func (s *FillOptionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFillOption(s)
	}
}

func (s *FillOptionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFillOption(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) FillOption() (localctx IFillOptionContext) {
	localctx = NewFillOptionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, SQLParserRULE_fillOption)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(528)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserT_NULL || _la == SQLParserT_PREVIOUS || _la == SQLParserL_INT || _la == SQLParserL_DEC) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IOrderByClauseContext is an interface to support dynamic dispatch.
type IOrderByClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOrderByClauseContext differentiates from other interfaces.
	IsOrderByClauseContext()
}

type OrderByClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrderByClauseContext() *OrderByClauseContext {
	var p = new(OrderByClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_orderByClause
	return p
}

func (*OrderByClauseContext) IsOrderByClauseContext() {}

func NewOrderByClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrderByClauseContext {
	var p = new(OrderByClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_orderByClause

	return p
}

func (s *OrderByClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *OrderByClauseContext) T_ORDER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ORDER, 0)
}

func (s *OrderByClauseContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *OrderByClauseContext) SortFields() ISortFieldsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISortFieldsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISortFieldsContext)
}

func (s *OrderByClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrderByClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OrderByClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterOrderByClause(s)
	}
}

func (s *OrderByClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitOrderByClause(s)
	}
}

func (s *OrderByClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitOrderByClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) OrderByClause() (localctx IOrderByClauseContext) {
	localctx = NewOrderByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, SQLParserRULE_orderByClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(530)
		p.Match(SQLParserT_ORDER)
	}
	{
		p.SetState(531)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(532)
		p.SortFields()
	}

	return localctx
}

// IIntervalByClauseContext is an interface to support dynamic dispatch.
type IIntervalByClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntervalByClauseContext differentiates from other interfaces.
	IsIntervalByClauseContext()
}

type IntervalByClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntervalByClauseContext() *IntervalByClauseContext {
	var p = new(IntervalByClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_intervalByClause
	return p
}

func (*IntervalByClauseContext) IsIntervalByClauseContext() {}

func NewIntervalByClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntervalByClauseContext {
	var p = new(IntervalByClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_intervalByClause

	return p
}

func (s *IntervalByClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *IntervalByClauseContext) T_INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL, 0)
}

func (s *IntervalByClauseContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *IntervalByClauseContext) IntervalNameVal() IIntervalNameValContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntervalNameValContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntervalNameValContext)
}

func (s *IntervalByClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntervalByClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntervalByClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIntervalByClause(s)
	}
}

func (s *IntervalByClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIntervalByClause(s)
	}
}

func (s *IntervalByClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIntervalByClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) IntervalByClause() (localctx IIntervalByClauseContext) {
	localctx = NewIntervalByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, SQLParserRULE_intervalByClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(534)
		p.Match(SQLParserT_INTERVAL)
	}
	{
		p.SetState(535)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(536)
		p.IntervalNameVal()
	}

	return localctx
}

// ISortFieldContext is an interface to support dynamic dispatch.
type ISortFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSortFieldContext differentiates from other interfaces.
	IsSortFieldContext()
}

type SortFieldContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortFieldContext() *SortFieldContext {
	var p = new(SortFieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_sortField
	return p
}

func (*SortFieldContext) IsSortFieldContext() {}

func NewSortFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortFieldContext {
	var p = new(SortFieldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_sortField

	return p
}

func (s *SortFieldContext) GetParser() antlr.Parser { return s.parser }

func (s *SortFieldContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *SortFieldContext) AllT_ASC() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_ASC)
}

func (s *SortFieldContext) T_ASC(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_ASC, i)
}

func (s *SortFieldContext) AllT_DESC() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_DESC)
}

func (s *SortFieldContext) T_DESC(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_DESC, i)
}

func (s *SortFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSortField(s)
	}
}

func (s *SortFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSortField(s)
	}
}

func (s *SortFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitSortField(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) SortField() (localctx ISortFieldContext) {
	localctx = NewSortFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, SQLParserRULE_sortField)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(538)
		p.expr(0)
	}
	p.SetState(542)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserT_ASC || _la == SQLParserT_DESC {
		{
			p.SetState(539)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_ASC || _la == SQLParserT_DESC) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(544)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ISortFieldsContext is an interface to support dynamic dispatch.
type ISortFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSortFieldsContext differentiates from other interfaces.
	IsSortFieldsContext()
}

type SortFieldsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortFieldsContext() *SortFieldsContext {
	var p = new(SortFieldsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_sortFields
	return p
}

func (*SortFieldsContext) IsSortFieldsContext() {}

func NewSortFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortFieldsContext {
	var p = new(SortFieldsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_sortFields

	return p
}

func (s *SortFieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *SortFieldsContext) AllSortField() []ISortFieldContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISortFieldContext)(nil)).Elem())
	var tst = make([]ISortFieldContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISortFieldContext)
		}
	}

	return tst
}

func (s *SortFieldsContext) SortField(i int) ISortFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISortFieldContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISortFieldContext)
}

func (s *SortFieldsContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *SortFieldsContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *SortFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSortFields(s)
	}
}

func (s *SortFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSortFields(s)
	}
}

func (s *SortFieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitSortFields(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) SortFields() (localctx ISortFieldsContext) {
	localctx = NewSortFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, SQLParserRULE_sortFields)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(545)
		p.SortField()
	}
	p.SetState(550)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserT_COMMA {
		{
			p.SetState(546)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(547)
			p.SortField()
		}

		p.SetState(552)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IHavingClauseContext is an interface to support dynamic dispatch.
type IHavingClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHavingClauseContext differentiates from other interfaces.
	IsHavingClauseContext()
}

type HavingClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHavingClauseContext() *HavingClauseContext {
	var p = new(HavingClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_havingClause
	return p
}

func (*HavingClauseContext) IsHavingClauseContext() {}

func NewHavingClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HavingClauseContext {
	var p = new(HavingClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_havingClause

	return p
}

func (s *HavingClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *HavingClauseContext) T_HAVING() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HAVING, 0)
}

func (s *HavingClauseContext) BoolExpr() IBoolExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBoolExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBoolExprContext)
}

func (s *HavingClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HavingClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HavingClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterHavingClause(s)
	}
}

func (s *HavingClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitHavingClause(s)
	}
}

func (s *HavingClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitHavingClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) HavingClause() (localctx IHavingClauseContext) {
	localctx = NewHavingClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, SQLParserRULE_havingClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(553)
		p.Match(SQLParserT_HAVING)
	}
	{
		p.SetState(554)
		p.boolExpr(0)
	}

	return localctx
}

// IBoolExprContext is an interface to support dynamic dispatch.
type IBoolExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBoolExprContext differentiates from other interfaces.
	IsBoolExprContext()
}

type BoolExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBoolExprContext() *BoolExprContext {
	var p = new(BoolExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_boolExpr
	return p
}

func (*BoolExprContext) IsBoolExprContext() {}

func NewBoolExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoolExprContext {
	var p = new(BoolExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_boolExpr

	return p
}

func (s *BoolExprContext) GetParser() antlr.Parser { return s.parser }

func (s *BoolExprContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *BoolExprContext) AllBoolExpr() []IBoolExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IBoolExprContext)(nil)).Elem())
	var tst = make([]IBoolExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IBoolExprContext)
		}
	}

	return tst
}

func (s *BoolExprContext) BoolExpr(i int) IBoolExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBoolExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IBoolExprContext)
}

func (s *BoolExprContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *BoolExprContext) BoolExprAtom() IBoolExprAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBoolExprAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBoolExprAtomContext)
}

func (s *BoolExprContext) BoolExprLogicalOp() IBoolExprLogicalOpContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBoolExprLogicalOpContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBoolExprLogicalOpContext)
}

func (s *BoolExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BoolExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBoolExpr(s)
	}
}

func (s *BoolExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBoolExpr(s)
	}
}

func (s *BoolExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBoolExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) BoolExpr() (localctx IBoolExprContext) {
	return p.boolExpr(0)
}

func (p *SQLParser) boolExpr(_p int) (localctx IBoolExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewBoolExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IBoolExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 118
	p.EnterRecursionRule(localctx, 118, SQLParserRULE_boolExpr, _p)

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(562)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(557)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(558)
			p.boolExpr(0)
		}
		{
			p.SetState(559)
			p.Match(SQLParserT_CLOSE_P)
		}

	case 2:
		{
			p.SetState(561)
			p.BoolExprAtom()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(570)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewBoolExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_boolExpr)
			p.SetState(564)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(565)
				p.BoolExprLogicalOp()
			}
			{
				p.SetState(566)
				p.boolExpr(3)
			}

		}
		p.SetState(572)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext())
	}

	return localctx
}

// IBoolExprLogicalOpContext is an interface to support dynamic dispatch.
type IBoolExprLogicalOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBoolExprLogicalOpContext differentiates from other interfaces.
	IsBoolExprLogicalOpContext()
}

type BoolExprLogicalOpContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBoolExprLogicalOpContext() *BoolExprLogicalOpContext {
	var p = new(BoolExprLogicalOpContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_boolExprLogicalOp
	return p
}

func (*BoolExprLogicalOpContext) IsBoolExprLogicalOpContext() {}

func NewBoolExprLogicalOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoolExprLogicalOpContext {
	var p = new(BoolExprLogicalOpContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_boolExprLogicalOp

	return p
}

func (s *BoolExprLogicalOpContext) GetParser() antlr.Parser { return s.parser }

func (s *BoolExprLogicalOpContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *BoolExprLogicalOpContext) T_OR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OR, 0)
}

func (s *BoolExprLogicalOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolExprLogicalOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BoolExprLogicalOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBoolExprLogicalOp(s)
	}
}

func (s *BoolExprLogicalOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBoolExprLogicalOp(s)
	}
}

func (s *BoolExprLogicalOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBoolExprLogicalOp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) BoolExprLogicalOp() (localctx IBoolExprLogicalOpContext) {
	localctx = NewBoolExprLogicalOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 120, SQLParserRULE_boolExprLogicalOp)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(573)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserT_AND || _la == SQLParserT_OR) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IBoolExprAtomContext is an interface to support dynamic dispatch.
type IBoolExprAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBoolExprAtomContext differentiates from other interfaces.
	IsBoolExprAtomContext()
}

type BoolExprAtomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBoolExprAtomContext() *BoolExprAtomContext {
	var p = new(BoolExprAtomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_boolExprAtom
	return p
}

func (*BoolExprAtomContext) IsBoolExprAtomContext() {}

func NewBoolExprAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoolExprAtomContext {
	var p = new(BoolExprAtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_boolExprAtom

	return p
}

func (s *BoolExprAtomContext) GetParser() antlr.Parser { return s.parser }

func (s *BoolExprAtomContext) BoolExprBinary() IBoolExprBinaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBoolExprBinaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBoolExprBinaryContext)
}

func (s *BoolExprAtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolExprAtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BoolExprAtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBoolExprAtom(s)
	}
}

func (s *BoolExprAtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBoolExprAtom(s)
	}
}

func (s *BoolExprAtomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBoolExprAtom(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) BoolExprAtom() (localctx IBoolExprAtomContext) {
	localctx = NewBoolExprAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, SQLParserRULE_boolExprAtom)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(575)
		p.BoolExprBinary()
	}

	return localctx
}

// IBoolExprBinaryContext is an interface to support dynamic dispatch.
type IBoolExprBinaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBoolExprBinaryContext differentiates from other interfaces.
	IsBoolExprBinaryContext()
}

type BoolExprBinaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBoolExprBinaryContext() *BoolExprBinaryContext {
	var p = new(BoolExprBinaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_boolExprBinary
	return p
}

func (*BoolExprBinaryContext) IsBoolExprBinaryContext() {}

func NewBoolExprBinaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoolExprBinaryContext {
	var p = new(BoolExprBinaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_boolExprBinary

	return p
}

func (s *BoolExprBinaryContext) GetParser() antlr.Parser { return s.parser }

func (s *BoolExprBinaryContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *BoolExprBinaryContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *BoolExprBinaryContext) BoolExprBinaryOperator() IBoolExprBinaryOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBoolExprBinaryOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBoolExprBinaryOperatorContext)
}

func (s *BoolExprBinaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolExprBinaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BoolExprBinaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBoolExprBinary(s)
	}
}

func (s *BoolExprBinaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBoolExprBinary(s)
	}
}

func (s *BoolExprBinaryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBoolExprBinary(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) BoolExprBinary() (localctx IBoolExprBinaryContext) {
	localctx = NewBoolExprBinaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, SQLParserRULE_boolExprBinary)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(577)
		p.expr(0)
	}
	{
		p.SetState(578)
		p.BoolExprBinaryOperator()
	}
	{
		p.SetState(579)
		p.expr(0)
	}

	return localctx
}

// IBoolExprBinaryOperatorContext is an interface to support dynamic dispatch.
type IBoolExprBinaryOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBoolExprBinaryOperatorContext differentiates from other interfaces.
	IsBoolExprBinaryOperatorContext()
}

type BoolExprBinaryOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBoolExprBinaryOperatorContext() *BoolExprBinaryOperatorContext {
	var p = new(BoolExprBinaryOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_boolExprBinaryOperator
	return p
}

func (*BoolExprBinaryOperatorContext) IsBoolExprBinaryOperatorContext() {}

func NewBoolExprBinaryOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoolExprBinaryOperatorContext {
	var p = new(BoolExprBinaryOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_boolExprBinaryOperator

	return p
}

func (s *BoolExprBinaryOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *BoolExprBinaryOperatorContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *BoolExprBinaryOperatorContext) T_NOTEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL, 0)
}

func (s *BoolExprBinaryOperatorContext) T_NOTEQUAL2() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL2, 0)
}

func (s *BoolExprBinaryOperatorContext) T_LESS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LESS, 0)
}

func (s *BoolExprBinaryOperatorContext) T_LESSEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LESSEQUAL, 0)
}

func (s *BoolExprBinaryOperatorContext) T_GREATER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GREATER, 0)
}

func (s *BoolExprBinaryOperatorContext) T_GREATEREQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GREATEREQUAL, 0)
}

func (s *BoolExprBinaryOperatorContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *BoolExprBinaryOperatorContext) T_REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REGEXP, 0)
}

func (s *BoolExprBinaryOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolExprBinaryOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BoolExprBinaryOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBoolExprBinaryOperator(s)
	}
}

func (s *BoolExprBinaryOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBoolExprBinaryOperator(s)
	}
}

func (s *BoolExprBinaryOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBoolExprBinaryOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) BoolExprBinaryOperator() (localctx IBoolExprBinaryOperatorContext) {
	localctx = NewBoolExprBinaryOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, SQLParserRULE_boolExprBinaryOperator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(589)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_EQUAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(581)
			p.Match(SQLParserT_EQUAL)
		}

	case SQLParserT_NOTEQUAL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(582)
			p.Match(SQLParserT_NOTEQUAL)
		}

	case SQLParserT_NOTEQUAL2:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(583)
			p.Match(SQLParserT_NOTEQUAL2)
		}

	case SQLParserT_LESS:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(584)
			p.Match(SQLParserT_LESS)
		}

	case SQLParserT_LESSEQUAL:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(585)
			p.Match(SQLParserT_LESSEQUAL)
		}

	case SQLParserT_GREATER:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(586)
			p.Match(SQLParserT_GREATER)
		}

	case SQLParserT_GREATEREQUAL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(587)
			p.Match(SQLParserT_GREATEREQUAL)
		}

	case SQLParserT_LIKE, SQLParserT_REGEXP:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(588)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_LIKE || _la == SQLParserT_REGEXP) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_expr
	return p
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *ExprContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *ExprContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *ExprContext) ExprFunc() IExprFuncContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprFuncContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprFuncContext)
}

func (s *ExprContext) ExprAtom() IExprAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprAtomContext)
}

func (s *ExprContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *ExprContext) T_MUL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MUL, 0)
}

func (s *ExprContext) T_DIV() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DIV, 0)
}

func (s *ExprContext) T_ADD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ADD, 0)
}

func (s *ExprContext) T_SUB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUB, 0)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterExpr(s)
	}
}

func (s *ExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitExpr(s)
	}
}

func (s *ExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Expr() (localctx IExprContext) {
	return p.expr(0)
}

func (p *SQLParser) expr(_p int) (localctx IExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 128
	p.EnterRecursionRule(localctx, 128, SQLParserRULE_expr, _p)

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(599)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 49, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(592)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(593)
			p.expr(0)
		}
		{
			p.SetState(594)
			p.Match(SQLParserT_CLOSE_P)
		}

	case 2:
		{
			p.SetState(596)
			p.ExprFunc()
		}

	case 3:
		{
			p.SetState(597)
			p.ExprAtom()
		}

	case 4:
		{
			p.SetState(598)
			p.DurationLit()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(615)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(613)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 50, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_expr)
				p.SetState(601)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(602)
					p.Match(SQLParserT_MUL)
				}
				{
					p.SetState(603)
					p.expr(9)
				}

			case 2:
				localctx = NewExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_expr)
				p.SetState(604)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(605)
					p.Match(SQLParserT_DIV)
				}
				{
					p.SetState(606)
					p.expr(8)
				}

			case 3:
				localctx = NewExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_expr)
				p.SetState(607)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				}
				{
					p.SetState(608)
					p.Match(SQLParserT_ADD)
				}
				{
					p.SetState(609)
					p.expr(7)
				}

			case 4:
				localctx = NewExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_expr)
				p.SetState(610)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				{
					p.SetState(611)
					p.Match(SQLParserT_SUB)
				}
				{
					p.SetState(612)
					p.expr(6)
				}

			}

		}
		p.SetState(617)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext())
	}

	return localctx
}

// IDurationLitContext is an interface to support dynamic dispatch.
type IDurationLitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDurationLitContext differentiates from other interfaces.
	IsDurationLitContext()
}

type DurationLitContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDurationLitContext() *DurationLitContext {
	var p = new(DurationLitContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_durationLit
	return p
}

func (*DurationLitContext) IsDurationLitContext() {}

func NewDurationLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DurationLitContext {
	var p = new(DurationLitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_durationLit

	return p
}

func (s *DurationLitContext) GetParser() antlr.Parser { return s.parser }

func (s *DurationLitContext) IntNumber() IIntNumberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntNumberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntNumberContext)
}

func (s *DurationLitContext) IntervalItem() IIntervalItemContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntervalItemContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntervalItemContext)
}

func (s *DurationLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DurationLitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DurationLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDurationLit(s)
	}
}

func (s *DurationLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDurationLit(s)
	}
}

func (s *DurationLitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDurationLit(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) DurationLit() (localctx IDurationLitContext) {
	localctx = NewDurationLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 130, SQLParserRULE_durationLit)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(618)
		p.IntNumber()
	}
	{
		p.SetState(619)
		p.IntervalItem()
	}

	return localctx
}

// IIntervalItemContext is an interface to support dynamic dispatch.
type IIntervalItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntervalItemContext differentiates from other interfaces.
	IsIntervalItemContext()
}

type IntervalItemContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntervalItemContext() *IntervalItemContext {
	var p = new(IntervalItemContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_intervalItem
	return p
}

func (*IntervalItemContext) IsIntervalItemContext() {}

func NewIntervalItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntervalItemContext {
	var p = new(IntervalItemContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_intervalItem

	return p
}

func (s *IntervalItemContext) GetParser() antlr.Parser { return s.parser }

func (s *IntervalItemContext) T_SECOND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SECOND, 0)
}

func (s *IntervalItemContext) T_MINUTE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MINUTE, 0)
}

func (s *IntervalItemContext) T_HOUR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HOUR, 0)
}

func (s *IntervalItemContext) T_DAY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DAY, 0)
}

func (s *IntervalItemContext) T_WEEK() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WEEK, 0)
}

func (s *IntervalItemContext) T_MONTH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MONTH, 0)
}

func (s *IntervalItemContext) T_YEAR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_YEAR, 0)
}

func (s *IntervalItemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntervalItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntervalItemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIntervalItem(s)
	}
}

func (s *IntervalItemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIntervalItem(s)
	}
}

func (s *IntervalItemContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIntervalItem(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) IntervalItem() (localctx IIntervalItemContext) {
	localctx = NewIntervalItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 132, SQLParserRULE_intervalItem)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(621)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-60)&-(0x1f+1)) == 0 && ((1<<uint((_la-60)))&((1<<(SQLParserT_SECOND-60))|(1<<(SQLParserT_MINUTE-60))|(1<<(SQLParserT_HOUR-60))|(1<<(SQLParserT_DAY-60))|(1<<(SQLParserT_WEEK-60))|(1<<(SQLParserT_MONTH-60))|(1<<(SQLParserT_YEAR-60)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IExprFuncContext is an interface to support dynamic dispatch.
type IExprFuncContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprFuncContext differentiates from other interfaces.
	IsExprFuncContext()
}

type ExprFuncContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprFuncContext() *ExprFuncContext {
	var p = new(ExprFuncContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_exprFunc
	return p
}

func (*ExprFuncContext) IsExprFuncContext() {}

func NewExprFuncContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprFuncContext {
	var p = new(ExprFuncContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_exprFunc

	return p
}

func (s *ExprFuncContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprFuncContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ExprFuncContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *ExprFuncContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *ExprFuncContext) ExprFuncParams() IExprFuncParamsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprFuncParamsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprFuncParamsContext)
}

func (s *ExprFuncContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprFuncContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprFuncContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterExprFunc(s)
	}
}

func (s *ExprFuncContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitExprFunc(s)
	}
}

func (s *ExprFuncContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitExprFunc(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ExprFunc() (localctx IExprFuncContext) {
	localctx = NewExprFuncContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 134, SQLParserRULE_exprFunc)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(623)
		p.Ident()
	}
	{
		p.SetState(624)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(626)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SQLParserT_CREATE)|(1<<SQLParserT_INTERVAL)|(1<<SQLParserT_SHARD)|(1<<SQLParserT_REPLICATION)|(1<<SQLParserT_TTL)|(1<<SQLParserT_KILL)|(1<<SQLParserT_ON)|(1<<SQLParserT_SHOW)|(1<<SQLParserT_DATASBAE)|(1<<SQLParserT_DATASBAES)|(1<<SQLParserT_NODE)|(1<<SQLParserT_MEASUREMENTS)|(1<<SQLParserT_MEASUREMENT)|(1<<SQLParserT_FIELD)|(1<<SQLParserT_TAG)|(1<<SQLParserT_KEYS)|(1<<SQLParserT_KEY)|(1<<SQLParserT_WITH)|(1<<SQLParserT_VALUES)|(1<<SQLParserT_FROM)|(1<<SQLParserT_WHERE)|(1<<SQLParserT_LIMIT))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(SQLParserT_QUERIES-32))|(1<<(SQLParserT_QUERY-32))|(1<<(SQLParserT_SELECT-32))|(1<<(SQLParserT_AS-32))|(1<<(SQLParserT_AND-32))|(1<<(SQLParserT_OR-32))|(1<<(SQLParserT_FILL-32))|(1<<(SQLParserT_NULL-32))|(1<<(SQLParserT_PREVIOUS-32))|(1<<(SQLParserT_ORDER-32))|(1<<(SQLParserT_ASC-32))|(1<<(SQLParserT_DESC-32))|(1<<(SQLParserT_LIKE-32))|(1<<(SQLParserT_NOT-32))|(1<<(SQLParserT_BETWEEN-32))|(1<<(SQLParserT_IS-32))|(1<<(SQLParserT_GROUP-32))|(1<<(SQLParserT_BY-32))|(1<<(SQLParserT_FOR-32))|(1<<(SQLParserT_STATS-32))|(1<<(SQLParserT_TIME-32))|(1<<(SQLParserT_PROFILE-32))|(1<<(SQLParserT_SECOND-32))|(1<<(SQLParserT_MINUTE-32))|(1<<(SQLParserT_HOUR-32))|(1<<(SQLParserT_DAY-32)))) != 0) || (((_la-64)&-(0x1f+1)) == 0 && ((1<<uint((_la-64)))&((1<<(SQLParserT_WEEK-64))|(1<<(SQLParserT_MONTH-64))|(1<<(SQLParserT_YEAR-64))|(1<<(SQLParserT_OPEN_P-64))|(1<<(SQLParserT_ADD-64))|(1<<(SQLParserT_SUB-64))|(1<<(SQLParserL_ID-64))|(1<<(SQLParserL_INT-64))|(1<<(SQLParserL_DEC-64)))) != 0) {
		{
			p.SetState(625)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(628)
		p.Match(SQLParserT_CLOSE_P)
	}

	return localctx
}

// IExprFuncParamsContext is an interface to support dynamic dispatch.
type IExprFuncParamsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprFuncParamsContext differentiates from other interfaces.
	IsExprFuncParamsContext()
}

type ExprFuncParamsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprFuncParamsContext() *ExprFuncParamsContext {
	var p = new(ExprFuncParamsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_exprFuncParams
	return p
}

func (*ExprFuncParamsContext) IsExprFuncParamsContext() {}

func NewExprFuncParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprFuncParamsContext {
	var p = new(ExprFuncParamsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_exprFuncParams

	return p
}

func (s *ExprFuncParamsContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprFuncParamsContext) AllFuncParam() []IFuncParamContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFuncParamContext)(nil)).Elem())
	var tst = make([]IFuncParamContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFuncParamContext)
		}
	}

	return tst
}

func (s *ExprFuncParamsContext) FuncParam(i int) IFuncParamContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFuncParamContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFuncParamContext)
}

func (s *ExprFuncParamsContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *ExprFuncParamsContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *ExprFuncParamsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprFuncParamsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprFuncParamsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterExprFuncParams(s)
	}
}

func (s *ExprFuncParamsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitExprFuncParams(s)
	}
}

func (s *ExprFuncParamsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitExprFuncParams(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ExprFuncParams() (localctx IExprFuncParamsContext) {
	localctx = NewExprFuncParamsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 136, SQLParserRULE_exprFuncParams)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(630)
		p.FuncParam()
	}
	p.SetState(635)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserT_COMMA {
		{
			p.SetState(631)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(632)
			p.FuncParam()
		}

		p.SetState(637)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IFuncParamContext is an interface to support dynamic dispatch.
type IFuncParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFuncParamContext differentiates from other interfaces.
	IsFuncParamContext()
}

type FuncParamContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuncParamContext() *FuncParamContext {
	var p = new(FuncParamContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_funcParam
	return p
}

func (*FuncParamContext) IsFuncParamContext() {}

func NewFuncParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncParamContext {
	var p = new(FuncParamContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_funcParam

	return p
}

func (s *FuncParamContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncParamContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *FuncParamContext) TagBooleanExpr() ITagBooleanExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagBooleanExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagBooleanExprContext)
}

func (s *FuncParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FuncParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFuncParam(s)
	}
}

func (s *FuncParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFuncParam(s)
	}
}

func (s *FuncParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFuncParam(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) FuncParam() (localctx IFuncParamContext) {
	localctx = NewFuncParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 138, SQLParserRULE_funcParam)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(640)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 54, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(638)
			p.expr(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(639)
			p.tagBooleanExpr(0)
		}

	}

	return localctx
}

// IExprAtomContext is an interface to support dynamic dispatch.
type IExprAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprAtomContext differentiates from other interfaces.
	IsExprAtomContext()
}

type ExprAtomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprAtomContext() *ExprAtomContext {
	var p = new(ExprAtomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_exprAtom
	return p
}

func (*ExprAtomContext) IsExprAtomContext() {}

func NewExprAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprAtomContext {
	var p = new(ExprAtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_exprAtom

	return p
}

func (s *ExprAtomContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprAtomContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ExprAtomContext) IdentFilter() IIdentFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentFilterContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentFilterContext)
}

func (s *ExprAtomContext) DecNumber() IDecNumberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDecNumberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDecNumberContext)
}

func (s *ExprAtomContext) IntNumber() IIntNumberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntNumberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntNumberContext)
}

func (s *ExprAtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprAtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprAtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterExprAtom(s)
	}
}

func (s *ExprAtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitExprAtom(s)
	}
}

func (s *ExprAtomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitExprAtom(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ExprAtom() (localctx IExprAtomContext) {
	localctx = NewExprAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 140, SQLParserRULE_exprAtom)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(648)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 56, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(642)
			p.Ident()
		}
		p.SetState(644)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 55, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(643)
				p.IdentFilter()
			}

		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(646)
			p.DecNumber()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(647)
			p.IntNumber()
		}

	}

	return localctx
}

// IIdentFilterContext is an interface to support dynamic dispatch.
type IIdentFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIdentFilterContext differentiates from other interfaces.
	IsIdentFilterContext()
}

type IdentFilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentFilterContext() *IdentFilterContext {
	var p = new(IdentFilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_identFilter
	return p
}

func (*IdentFilterContext) IsIdentFilterContext() {}

func NewIdentFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentFilterContext {
	var p = new(IdentFilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_identFilter

	return p
}

func (s *IdentFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentFilterContext) T_OPEN_SB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_SB, 0)
}

func (s *IdentFilterContext) TagBooleanExpr() ITagBooleanExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagBooleanExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagBooleanExprContext)
}

func (s *IdentFilterContext) T_CLOSE_SB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_SB, 0)
}

func (s *IdentFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIdentFilter(s)
	}
}

func (s *IdentFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIdentFilter(s)
	}
}

func (s *IdentFilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIdentFilter(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) IdentFilter() (localctx IIdentFilterContext) {
	localctx = NewIdentFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 142, SQLParserRULE_identFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(650)
		p.Match(SQLParserT_OPEN_SB)
	}
	{
		p.SetState(651)
		p.tagBooleanExpr(0)
	}
	{
		p.SetState(652)
		p.Match(SQLParserT_CLOSE_SB)
	}

	return localctx
}

// IIntNumberContext is an interface to support dynamic dispatch.
type IIntNumberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntNumberContext differentiates from other interfaces.
	IsIntNumberContext()
}

type IntNumberContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntNumberContext() *IntNumberContext {
	var p = new(IntNumberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_intNumber
	return p
}

func (*IntNumberContext) IsIntNumberContext() {}

func NewIntNumberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntNumberContext {
	var p = new(IntNumberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_intNumber

	return p
}

func (s *IntNumberContext) GetParser() antlr.Parser { return s.parser }

func (s *IntNumberContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *IntNumberContext) T_SUB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUB, 0)
}

func (s *IntNumberContext) T_ADD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ADD, 0)
}

func (s *IntNumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntNumberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntNumberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIntNumber(s)
	}
}

func (s *IntNumberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIntNumber(s)
	}
}

func (s *IntNumberContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIntNumber(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) IntNumber() (localctx IIntNumberContext) {
	localctx = NewIntNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 144, SQLParserRULE_intNumber)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(655)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(654)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_ADD || _la == SQLParserT_SUB) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(657)
		p.Match(SQLParserL_INT)
	}

	return localctx
}

// IDecNumberContext is an interface to support dynamic dispatch.
type IDecNumberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDecNumberContext differentiates from other interfaces.
	IsDecNumberContext()
}

type DecNumberContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDecNumberContext() *DecNumberContext {
	var p = new(DecNumberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_decNumber
	return p
}

func (*DecNumberContext) IsDecNumberContext() {}

func NewDecNumberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DecNumberContext {
	var p = new(DecNumberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_decNumber

	return p
}

func (s *DecNumberContext) GetParser() antlr.Parser { return s.parser }

func (s *DecNumberContext) L_DEC() antlr.TerminalNode {
	return s.GetToken(SQLParserL_DEC, 0)
}

func (s *DecNumberContext) T_SUB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUB, 0)
}

func (s *DecNumberContext) T_ADD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ADD, 0)
}

func (s *DecNumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DecNumberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DecNumberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDecNumber(s)
	}
}

func (s *DecNumberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDecNumber(s)
	}
}

func (s *DecNumberContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDecNumber(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) DecNumber() (localctx IDecNumberContext) {
	localctx = NewDecNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 146, SQLParserRULE_decNumber)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(660)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(659)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_ADD || _la == SQLParserT_SUB) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(662)
		p.Match(SQLParserL_DEC)
	}

	return localctx
}

// ILimitClauseContext is an interface to support dynamic dispatch.
type ILimitClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLimitClauseContext differentiates from other interfaces.
	IsLimitClauseContext()
}

type LimitClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLimitClauseContext() *LimitClauseContext {
	var p = new(LimitClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_limitClause
	return p
}

func (*LimitClauseContext) IsLimitClauseContext() {}

func NewLimitClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitClauseContext {
	var p = new(LimitClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_limitClause

	return p
}

func (s *LimitClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitClauseContext) T_LIMIT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIMIT, 0)
}

func (s *LimitClauseContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *LimitClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LimitClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterLimitClause(s)
	}
}

func (s *LimitClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitLimitClause(s)
	}
}

func (s *LimitClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitLimitClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) LimitClause() (localctx ILimitClauseContext) {
	localctx = NewLimitClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 148, SQLParserRULE_limitClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(664)
		p.Match(SQLParserT_LIMIT)
	}
	{
		p.SetState(665)
		p.Match(SQLParserL_INT)
	}

	return localctx
}

// IMetricNameContext is an interface to support dynamic dispatch.
type IMetricNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMetricNameContext differentiates from other interfaces.
	IsMetricNameContext()
}

type MetricNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetricNameContext() *MetricNameContext {
	var p = new(MetricNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_metricName
	return p
}

func (*MetricNameContext) IsMetricNameContext() {}

func NewMetricNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MetricNameContext {
	var p = new(MetricNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_metricName

	return p
}

func (s *MetricNameContext) GetParser() antlr.Parser { return s.parser }

func (s *MetricNameContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *MetricNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MetricNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MetricNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterMetricName(s)
	}
}

func (s *MetricNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitMetricName(s)
	}
}

func (s *MetricNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitMetricName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) MetricName() (localctx IMetricNameContext) {
	localctx = NewMetricNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 150, SQLParserRULE_metricName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(667)
		p.Ident()
	}

	return localctx
}

// ITagKeyContext is an interface to support dynamic dispatch.
type ITagKeyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagKeyContext differentiates from other interfaces.
	IsTagKeyContext()
}

type TagKeyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagKeyContext() *TagKeyContext {
	var p = new(TagKeyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagKey
	return p
}

func (*TagKeyContext) IsTagKeyContext() {}

func NewTagKeyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagKeyContext {
	var p = new(TagKeyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagKey

	return p
}

func (s *TagKeyContext) GetParser() antlr.Parser { return s.parser }

func (s *TagKeyContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *TagKeyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagKeyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TagKeyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagKey(s)
	}
}

func (s *TagKeyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagKey(s)
	}
}

func (s *TagKeyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagKey(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TagKey() (localctx ITagKeyContext) {
	localctx = NewTagKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 152, SQLParserRULE_tagKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(669)
		p.Ident()
	}

	return localctx
}

// ITagValueContext is an interface to support dynamic dispatch.
type ITagValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagValueContext differentiates from other interfaces.
	IsTagValueContext()
}

type TagValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagValueContext() *TagValueContext {
	var p = new(TagValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagValue
	return p
}

func (*TagValueContext) IsTagValueContext() {}

func NewTagValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagValueContext {
	var p = new(TagValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagValue

	return p
}

func (s *TagValueContext) GetParser() antlr.Parser { return s.parser }

func (s *TagValueContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *TagValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TagValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagValue(s)
	}
}

func (s *TagValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagValue(s)
	}
}

func (s *TagValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TagValue() (localctx ITagValueContext) {
	localctx = NewTagValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 154, SQLParserRULE_tagValue)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(671)
		p.Ident()
	}

	return localctx
}

// ITagValuePatternContext is an interface to support dynamic dispatch.
type ITagValuePatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagValuePatternContext differentiates from other interfaces.
	IsTagValuePatternContext()
}

type TagValuePatternContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagValuePatternContext() *TagValuePatternContext {
	var p = new(TagValuePatternContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagValuePattern
	return p
}

func (*TagValuePatternContext) IsTagValuePatternContext() {}

func NewTagValuePatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagValuePatternContext {
	var p = new(TagValuePatternContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagValuePattern

	return p
}

func (s *TagValuePatternContext) GetParser() antlr.Parser { return s.parser }

func (s *TagValuePatternContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *TagValuePatternContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagValuePatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TagValuePatternContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagValuePattern(s)
	}
}

func (s *TagValuePatternContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagValuePattern(s)
	}
}

func (s *TagValuePatternContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagValuePattern(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) TagValuePattern() (localctx ITagValuePatternContext) {
	localctx = NewTagValuePatternContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 156, SQLParserRULE_tagValuePattern)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(673)
		p.Ident()
	}

	return localctx
}

// IIdentContext is an interface to support dynamic dispatch.
type IIdentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIdentContext differentiates from other interfaces.
	IsIdentContext()
}

type IdentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentContext() *IdentContext {
	var p = new(IdentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_ident
	return p
}

func (*IdentContext) IsIdentContext() {}

func NewIdentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentContext {
	var p = new(IdentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_ident

	return p
}

func (s *IdentContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentContext) AllL_ID() []antlr.TerminalNode {
	return s.GetTokens(SQLParserL_ID)
}

func (s *IdentContext) L_ID(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserL_ID, i)
}

func (s *IdentContext) AllNonReservedWords() []INonReservedWordsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*INonReservedWordsContext)(nil)).Elem())
	var tst = make([]INonReservedWordsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(INonReservedWordsContext)
		}
	}

	return tst
}

func (s *IdentContext) NonReservedWords(i int) INonReservedWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INonReservedWordsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(INonReservedWordsContext)
}

func (s *IdentContext) AllT_DOT() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_DOT)
}

func (s *IdentContext) T_DOT(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_DOT, i)
}

func (s *IdentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIdent(s)
	}
}

func (s *IdentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIdent(s)
	}
}

func (s *IdentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIdent(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Ident() (localctx IIdentContext) {
	localctx = NewIdentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 158, SQLParserRULE_ident)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(677)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserL_ID:
		{
			p.SetState(675)
			p.Match(SQLParserL_ID)
		}

	case SQLParserT_CREATE, SQLParserT_INTERVAL, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NODE, SQLParserT_MEASUREMENTS, SQLParserT_MEASUREMENT, SQLParserT_FIELD, SQLParserT_TAG, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_PROFILE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
		{
			p.SetState(676)
			p.NonReservedWords()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(686)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 61, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(679)
				p.Match(SQLParserT_DOT)
			}
			p.SetState(682)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case SQLParserL_ID:
				{
					p.SetState(680)
					p.Match(SQLParserL_ID)
				}

			case SQLParserT_CREATE, SQLParserT_INTERVAL, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NODE, SQLParserT_MEASUREMENTS, SQLParserT_MEASUREMENT, SQLParserT_FIELD, SQLParserT_TAG, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_PROFILE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
				{
					p.SetState(681)
					p.NonReservedWords()
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

		}
		p.SetState(688)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 61, p.GetParserRuleContext())
	}

	return localctx
}

// INonReservedWordsContext is an interface to support dynamic dispatch.
type INonReservedWordsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNonReservedWordsContext differentiates from other interfaces.
	IsNonReservedWordsContext()
}

type NonReservedWordsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNonReservedWordsContext() *NonReservedWordsContext {
	var p = new(NonReservedWordsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_nonReservedWords
	return p
}

func (*NonReservedWordsContext) IsNonReservedWordsContext() {}

func NewNonReservedWordsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NonReservedWordsContext {
	var p = new(NonReservedWordsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_nonReservedWords

	return p
}

func (s *NonReservedWordsContext) GetParser() antlr.Parser { return s.parser }

func (s *NonReservedWordsContext) T_CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CREATE, 0)
}

func (s *NonReservedWordsContext) T_INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL, 0)
}

func (s *NonReservedWordsContext) T_SHARD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHARD, 0)
}

func (s *NonReservedWordsContext) T_REPLICATION() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REPLICATION, 0)
}

func (s *NonReservedWordsContext) T_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TTL, 0)
}

func (s *NonReservedWordsContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *NonReservedWordsContext) T_KILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KILL, 0)
}

func (s *NonReservedWordsContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *NonReservedWordsContext) T_DATASBAES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAES, 0)
}

func (s *NonReservedWordsContext) T_NODE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NODE, 0)
}

func (s *NonReservedWordsContext) T_MEASUREMENTS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MEASUREMENTS, 0)
}

func (s *NonReservedWordsContext) T_MEASUREMENT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MEASUREMENT, 0)
}

func (s *NonReservedWordsContext) T_FIELD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIELD, 0)
}

func (s *NonReservedWordsContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *NonReservedWordsContext) T_KEYS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEYS, 0)
}

func (s *NonReservedWordsContext) T_KEY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEY, 0)
}

func (s *NonReservedWordsContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *NonReservedWordsContext) T_VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUES, 0)
}

func (s *NonReservedWordsContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *NonReservedWordsContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *NonReservedWordsContext) T_LIMIT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIMIT, 0)
}

func (s *NonReservedWordsContext) T_QUERIES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUERIES, 0)
}

func (s *NonReservedWordsContext) T_QUERY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUERY, 0)
}

func (s *NonReservedWordsContext) T_SELECT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SELECT, 0)
}

func (s *NonReservedWordsContext) T_AS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AS, 0)
}

func (s *NonReservedWordsContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *NonReservedWordsContext) T_OR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OR, 0)
}

func (s *NonReservedWordsContext) T_NULL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NULL, 0)
}

func (s *NonReservedWordsContext) T_PREVIOUS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PREVIOUS, 0)
}

func (s *NonReservedWordsContext) T_FILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FILL, 0)
}

func (s *NonReservedWordsContext) T_ORDER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ORDER, 0)
}

func (s *NonReservedWordsContext) T_ASC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ASC, 0)
}

func (s *NonReservedWordsContext) T_DESC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DESC, 0)
}

func (s *NonReservedWordsContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *NonReservedWordsContext) T_NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOT, 0)
}

func (s *NonReservedWordsContext) T_BETWEEN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BETWEEN, 0)
}

func (s *NonReservedWordsContext) T_IS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IS, 0)
}

func (s *NonReservedWordsContext) T_PROFILE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PROFILE, 0)
}

func (s *NonReservedWordsContext) T_GROUP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GROUP, 0)
}

func (s *NonReservedWordsContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *NonReservedWordsContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *NonReservedWordsContext) T_STATS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATS, 0)
}

func (s *NonReservedWordsContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *NonReservedWordsContext) T_FOR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FOR, 0)
}

func (s *NonReservedWordsContext) T_SECOND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SECOND, 0)
}

func (s *NonReservedWordsContext) T_MINUTE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MINUTE, 0)
}

func (s *NonReservedWordsContext) T_HOUR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HOUR, 0)
}

func (s *NonReservedWordsContext) T_DAY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DAY, 0)
}

func (s *NonReservedWordsContext) T_WEEK() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WEEK, 0)
}

func (s *NonReservedWordsContext) T_MONTH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MONTH, 0)
}

func (s *NonReservedWordsContext) T_YEAR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_YEAR, 0)
}

func (s *NonReservedWordsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NonReservedWordsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NonReservedWordsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNonReservedWords(s)
	}
}

func (s *NonReservedWordsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNonReservedWords(s)
	}
}

func (s *NonReservedWordsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitNonReservedWords(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) NonReservedWords() (localctx INonReservedWordsContext) {
	localctx = NewNonReservedWordsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 160, SQLParserRULE_nonReservedWords)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(689)
		_la = p.GetTokenStream().LA(1)

		if !((((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SQLParserT_CREATE)|(1<<SQLParserT_INTERVAL)|(1<<SQLParserT_SHARD)|(1<<SQLParserT_REPLICATION)|(1<<SQLParserT_TTL)|(1<<SQLParserT_KILL)|(1<<SQLParserT_ON)|(1<<SQLParserT_SHOW)|(1<<SQLParserT_DATASBAE)|(1<<SQLParserT_DATASBAES)|(1<<SQLParserT_NODE)|(1<<SQLParserT_MEASUREMENTS)|(1<<SQLParserT_MEASUREMENT)|(1<<SQLParserT_FIELD)|(1<<SQLParserT_TAG)|(1<<SQLParserT_KEYS)|(1<<SQLParserT_KEY)|(1<<SQLParserT_WITH)|(1<<SQLParserT_VALUES)|(1<<SQLParserT_FROM)|(1<<SQLParserT_WHERE)|(1<<SQLParserT_LIMIT))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(SQLParserT_QUERIES-32))|(1<<(SQLParserT_QUERY-32))|(1<<(SQLParserT_SELECT-32))|(1<<(SQLParserT_AS-32))|(1<<(SQLParserT_AND-32))|(1<<(SQLParserT_OR-32))|(1<<(SQLParserT_FILL-32))|(1<<(SQLParserT_NULL-32))|(1<<(SQLParserT_PREVIOUS-32))|(1<<(SQLParserT_ORDER-32))|(1<<(SQLParserT_ASC-32))|(1<<(SQLParserT_DESC-32))|(1<<(SQLParserT_LIKE-32))|(1<<(SQLParserT_NOT-32))|(1<<(SQLParserT_BETWEEN-32))|(1<<(SQLParserT_IS-32))|(1<<(SQLParserT_GROUP-32))|(1<<(SQLParserT_BY-32))|(1<<(SQLParserT_FOR-32))|(1<<(SQLParserT_STATS-32))|(1<<(SQLParserT_TIME-32))|(1<<(SQLParserT_PROFILE-32))|(1<<(SQLParserT_SECOND-32))|(1<<(SQLParserT_MINUTE-32))|(1<<(SQLParserT_HOUR-32))|(1<<(SQLParserT_DAY-32)))) != 0) || (((_la-64)&-(0x1f+1)) == 0 && ((1<<uint((_la-64)))&((1<<(SQLParserT_WEEK-64))|(1<<(SQLParserT_MONTH-64))|(1<<(SQLParserT_YEAR-64)))) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

func (p *SQLParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 41:
		var t *ClauseBooleanExprContext = nil
		if localctx != nil {
			t = localctx.(*ClauseBooleanExprContext)
		}
		return p.ClauseBooleanExpr_Sempred(t, predIndex)

	case 44:
		var t *TagBooleanExprContext = nil
		if localctx != nil {
			t = localctx.(*TagBooleanExprContext)
		}
		return p.TagBooleanExpr_Sempred(t, predIndex)

	case 59:
		var t *BoolExprContext = nil
		if localctx != nil {
			t = localctx.(*BoolExprContext)
		}
		return p.BoolExpr_Sempred(t, predIndex)

	case 64:
		var t *ExprContext = nil
		if localctx != nil {
			t = localctx.(*ExprContext)
		}
		return p.Expr_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *SQLParser) ClauseBooleanExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) TagBooleanExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 1:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) BoolExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 3:
		return p.Precpred(p.GetParserRuleContext(), 8)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 7)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 6)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 5)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
