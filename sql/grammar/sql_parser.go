// Code generated from /Users/dupeng/Documents/gohub/src/github.com/eleme/lindb/cmd/sql/antlr4/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser // SQL

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"reflect"
	"strconv"
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
	"statement", "statement_list", "create_database_stmt", "with_clause_list",
	"with_clause", "interval_define_list", "interval_define", "shard_num",
	"ttl_val", "metattl_val", "past_val", "future_val", "interval_name_val",
	"replica_factor", "database_name", "update_database_stmt", "drop_database_stmt",
	"show_databases_stmt", "show_node_stmt", "show_measurements_stmt", "show_tag_keys_stmt",
	"show_info_stmt", "show_tag_values_stmt", "show_tag_values_info_stmt",
	"show_field_keys_stmt", "show_queries_stmt", "show_stats_stmt", "with_measurement_clause",
	"with_tag_clause", "where_tag_cascade", "kill_query_stmt", "query_id",
	"server_id", "module", "component", "query_stmt", "fields", "field", "alias",
	"from_clause", "where_clause", "clause_boolean_expr", "tag_cascade_expr",
	"tag_equal_expr", "tag_boolean_expr", "tag_value_list", "time_expr", "time_boolean_expr",
	"now_expr", "now_func", "group_by_clause", "dimensions", "dimension", "fill_option",
	"order_by_clause", "interval_by_clause", "sort_field", "sort_fields", "having_clause",
	"bool_expr", "bool_expr_logical_op", "bool_expr_atom", "bool_expr_binary",
	"bool_expr_binary_operator", "expr", "duration_lit", "interval_item", "expr_func",
	"expr_func_params", "func_param", "expr_atom", "ident_filter", "int_number",
	"dec_number", "limit_clause", "metric_name", "tag_key", "tag_value", "tag_value_pattern",
	"ident", "non_reserved_words",
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
	SQLParserRULE_statement                 = 0
	SQLParserRULE_statement_list            = 1
	SQLParserRULE_create_database_stmt      = 2
	SQLParserRULE_with_clause_list          = 3
	SQLParserRULE_with_clause               = 4
	SQLParserRULE_interval_define_list      = 5
	SQLParserRULE_interval_define           = 6
	SQLParserRULE_shard_num                 = 7
	SQLParserRULE_ttl_val                   = 8
	SQLParserRULE_metattl_val               = 9
	SQLParserRULE_past_val                  = 10
	SQLParserRULE_future_val                = 11
	SQLParserRULE_interval_name_val         = 12
	SQLParserRULE_replica_factor            = 13
	SQLParserRULE_database_name             = 14
	SQLParserRULE_update_database_stmt      = 15
	SQLParserRULE_drop_database_stmt        = 16
	SQLParserRULE_show_databases_stmt       = 17
	SQLParserRULE_show_node_stmt            = 18
	SQLParserRULE_show_measurements_stmt    = 19
	SQLParserRULE_show_tag_keys_stmt        = 20
	SQLParserRULE_show_info_stmt            = 21
	SQLParserRULE_show_tag_values_stmt      = 22
	SQLParserRULE_show_tag_values_info_stmt = 23
	SQLParserRULE_show_field_keys_stmt      = 24
	SQLParserRULE_show_queries_stmt         = 25
	SQLParserRULE_show_stats_stmt           = 26
	SQLParserRULE_with_measurement_clause   = 27
	SQLParserRULE_with_tag_clause           = 28
	SQLParserRULE_where_tag_cascade         = 29
	SQLParserRULE_kill_query_stmt           = 30
	SQLParserRULE_query_id                  = 31
	SQLParserRULE_server_id                 = 32
	SQLParserRULE_module                    = 33
	SQLParserRULE_component                 = 34
	SQLParserRULE_query_stmt                = 35
	SQLParserRULE_fields                    = 36
	SQLParserRULE_field                     = 37
	SQLParserRULE_alias                     = 38
	SQLParserRULE_from_clause               = 39
	SQLParserRULE_where_clause              = 40
	SQLParserRULE_clause_boolean_expr       = 41
	SQLParserRULE_tag_cascade_expr          = 42
	SQLParserRULE_tag_equal_expr            = 43
	SQLParserRULE_tag_boolean_expr          = 44
	SQLParserRULE_tag_value_list            = 45
	SQLParserRULE_time_expr                 = 46
	SQLParserRULE_time_boolean_expr         = 47
	SQLParserRULE_now_expr                  = 48
	SQLParserRULE_now_func                  = 49
	SQLParserRULE_group_by_clause           = 50
	SQLParserRULE_dimensions                = 51
	SQLParserRULE_dimension                 = 52
	SQLParserRULE_fill_option               = 53
	SQLParserRULE_order_by_clause           = 54
	SQLParserRULE_interval_by_clause        = 55
	SQLParserRULE_sort_field                = 56
	SQLParserRULE_sort_fields               = 57
	SQLParserRULE_having_clause             = 58
	SQLParserRULE_bool_expr                 = 59
	SQLParserRULE_bool_expr_logical_op      = 60
	SQLParserRULE_bool_expr_atom            = 61
	SQLParserRULE_bool_expr_binary          = 62
	SQLParserRULE_bool_expr_binary_operator = 63
	SQLParserRULE_expr                      = 64
	SQLParserRULE_duration_lit              = 65
	SQLParserRULE_interval_item             = 66
	SQLParserRULE_expr_func                 = 67
	SQLParserRULE_expr_func_params          = 68
	SQLParserRULE_func_param                = 69
	SQLParserRULE_expr_atom                 = 70
	SQLParserRULE_ident_filter              = 71
	SQLParserRULE_int_number                = 72
	SQLParserRULE_dec_number                = 73
	SQLParserRULE_limit_clause              = 74
	SQLParserRULE_metric_name               = 75
	SQLParserRULE_tag_key                   = 76
	SQLParserRULE_tag_value                 = 77
	SQLParserRULE_tag_value_pattern         = 78
	SQLParserRULE_ident                     = 79
	SQLParserRULE_non_reserved_words        = 80
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

func (s *StatementContext) Statement_list() IStatement_listContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatement_listContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStatement_listContext)
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
		p.Statement_list()
	}
	{
		p.SetState(163)
		p.Match(SQLParserEOF)
	}

	return localctx
}

// IStatement_listContext is an interface to support dynamic dispatch.
type IStatement_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatement_listContext differentiates from other interfaces.
	IsStatement_listContext()
}

type Statement_listContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatement_listContext() *Statement_listContext {
	var p = new(Statement_listContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_statement_list
	return p
}

func (*Statement_listContext) IsStatement_listContext() {}

func NewStatement_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Statement_listContext {
	var p = new(Statement_listContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_statement_list

	return p
}

func (s *Statement_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Statement_listContext) Create_database_stmt() ICreate_database_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICreate_database_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICreate_database_stmtContext)
}

func (s *Statement_listContext) Update_database_stmt() IUpdate_database_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUpdate_database_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUpdate_database_stmtContext)
}

func (s *Statement_listContext) Drop_database_stmt() IDrop_database_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDrop_database_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDrop_database_stmtContext)
}

func (s *Statement_listContext) Show_databases_stmt() IShow_databases_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShow_databases_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShow_databases_stmtContext)
}

func (s *Statement_listContext) Show_node_stmt() IShow_node_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShow_node_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShow_node_stmtContext)
}

func (s *Statement_listContext) Show_measurements_stmt() IShow_measurements_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShow_measurements_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShow_measurements_stmtContext)
}

func (s *Statement_listContext) Show_info_stmt() IShow_info_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShow_info_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShow_info_stmtContext)
}

func (s *Statement_listContext) Show_tag_keys_stmt() IShow_tag_keys_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShow_tag_keys_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShow_tag_keys_stmtContext)
}

func (s *Statement_listContext) Show_queries_stmt() IShow_queries_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShow_queries_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShow_queries_stmtContext)
}

func (s *Statement_listContext) Show_tag_values_stmt() IShow_tag_values_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShow_tag_values_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShow_tag_values_stmtContext)
}

func (s *Statement_listContext) Show_tag_values_info_stmt() IShow_tag_values_info_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShow_tag_values_info_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShow_tag_values_info_stmtContext)
}

func (s *Statement_listContext) Show_field_keys_stmt() IShow_field_keys_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShow_field_keys_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShow_field_keys_stmtContext)
}

func (s *Statement_listContext) Show_stats_stmt() IShow_stats_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShow_stats_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShow_stats_stmtContext)
}

func (s *Statement_listContext) Kill_query_stmt() IKill_query_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IKill_query_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IKill_query_stmtContext)
}

func (s *Statement_listContext) Query_stmt() IQuery_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQuery_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IQuery_stmtContext)
}

func (s *Statement_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Statement_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Statement_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterStatement_list(s)
	}
}

func (s *Statement_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitStatement_list(s)
	}
}

func (s *Statement_listContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitStatement_list(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Statement_list() (localctx IStatement_listContext) {
	localctx = NewStatement_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, SQLParserRULE_statement_list)

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
			p.Create_database_stmt()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(166)
			p.Update_database_stmt()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(167)
			p.Drop_database_stmt()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(168)
			p.Show_databases_stmt()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(169)
			p.Show_node_stmt()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(170)
			p.Show_measurements_stmt()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(171)
			p.Show_info_stmt()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(172)
			p.Show_tag_keys_stmt()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(173)
			p.Show_queries_stmt()
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(174)
			p.Show_tag_values_stmt()
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(175)
			p.Show_tag_values_info_stmt()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(176)
			p.Show_field_keys_stmt()
		}

	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(177)
			p.Show_stats_stmt()
		}

	case 14:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(178)
			p.Kill_query_stmt()
		}

	case 15:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(179)
			p.Query_stmt()
		}

	}

	return localctx
}

// ICreate_database_stmtContext is an interface to support dynamic dispatch.
type ICreate_database_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCreate_database_stmtContext differentiates from other interfaces.
	IsCreate_database_stmtContext()
}

type Create_database_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreate_database_stmtContext() *Create_database_stmtContext {
	var p = new(Create_database_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_create_database_stmt
	return p
}

func (*Create_database_stmtContext) IsCreate_database_stmtContext() {}

func NewCreate_database_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Create_database_stmtContext {
	var p = new(Create_database_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_create_database_stmt

	return p
}

func (s *Create_database_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Create_database_stmtContext) T_CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CREATE, 0)
}

func (s *Create_database_stmtContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *Create_database_stmtContext) Database_name() IDatabase_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDatabase_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDatabase_nameContext)
}

func (s *Create_database_stmtContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *Create_database_stmtContext) With_clause_list() IWith_clause_listContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWith_clause_listContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWith_clause_listContext)
}

func (s *Create_database_stmtContext) T_COMMA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, 0)
}

func (s *Create_database_stmtContext) Interval_define_list() IInterval_define_listContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInterval_define_listContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInterval_define_listContext)
}

func (s *Create_database_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Create_database_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Create_database_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterCreate_database_stmt(s)
	}
}

func (s *Create_database_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitCreate_database_stmt(s)
	}
}

func (s *Create_database_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitCreate_database_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Create_database_stmt() (localctx ICreate_database_stmtContext) {
	localctx = NewCreate_database_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SQLParserRULE_create_database_stmt)
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
		p.Database_name()
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
			p.With_clause_list()
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
			p.Interval_define_list()
		}

	}

	return localctx
}

// IWith_clause_listContext is an interface to support dynamic dispatch.
type IWith_clause_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWith_clause_listContext differentiates from other interfaces.
	IsWith_clause_listContext()
}

type With_clause_listContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWith_clause_listContext() *With_clause_listContext {
	var p = new(With_clause_listContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_with_clause_list
	return p
}

func (*With_clause_listContext) IsWith_clause_listContext() {}

func NewWith_clause_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *With_clause_listContext {
	var p = new(With_clause_listContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_with_clause_list

	return p
}

func (s *With_clause_listContext) GetParser() antlr.Parser { return s.parser }

func (s *With_clause_listContext) AllWith_clause() []IWith_clauseContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWith_clauseContext)(nil)).Elem())
	var tst = make([]IWith_clauseContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWith_clauseContext)
		}
	}

	return tst
}

func (s *With_clause_listContext) With_clause(i int) IWith_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWith_clauseContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWith_clauseContext)
}

func (s *With_clause_listContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *With_clause_listContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *With_clause_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *With_clause_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *With_clause_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWith_clause_list(s)
	}
}

func (s *With_clause_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWith_clause_list(s)
	}
}

func (s *With_clause_listContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWith_clause_list(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) With_clause_list() (localctx IWith_clause_listContext) {
	localctx = NewWith_clause_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, SQLParserRULE_with_clause_list)

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
		p.With_clause()
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
				p.With_clause()
			}

		}
		p.SetState(200)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}

	return localctx
}

// IWith_clauseContext is an interface to support dynamic dispatch.
type IWith_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWith_clauseContext differentiates from other interfaces.
	IsWith_clauseContext()
}

type With_clauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWith_clauseContext() *With_clauseContext {
	var p = new(With_clauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_with_clause
	return p
}

func (*With_clauseContext) IsWith_clauseContext() {}

func NewWith_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *With_clauseContext {
	var p = new(With_clauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_with_clause

	return p
}

func (s *With_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *With_clauseContext) T_INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL, 0)
}

func (s *With_clauseContext) Duration_lit() IDuration_litContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDuration_litContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDuration_litContext)
}

func (s *With_clauseContext) T_SHARD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHARD, 0)
}

func (s *With_clauseContext) Shard_num() IShard_numContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShard_numContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShard_numContext)
}

func (s *With_clauseContext) T_REPLICATION() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REPLICATION, 0)
}

func (s *With_clauseContext) Replica_factor() IReplica_factorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IReplica_factorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IReplica_factorContext)
}

func (s *With_clauseContext) T_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TTL, 0)
}

func (s *With_clauseContext) Ttl_val() ITtl_valContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITtl_valContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITtl_valContext)
}

func (s *With_clauseContext) T_META_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_META_TTL, 0)
}

func (s *With_clauseContext) Metattl_val() IMetattl_valContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetattl_valContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetattl_valContext)
}

func (s *With_clauseContext) T_PAST_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PAST_TTL, 0)
}

func (s *With_clauseContext) Past_val() IPast_valContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPast_valContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPast_valContext)
}

func (s *With_clauseContext) T_FUTURE_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FUTURE_TTL, 0)
}

func (s *With_clauseContext) Future_val() IFuture_valContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFuture_valContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFuture_valContext)
}

func (s *With_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *With_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *With_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWith_clause(s)
	}
}

func (s *With_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWith_clause(s)
	}
}

func (s *With_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWith_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) With_clause() (localctx IWith_clauseContext) {
	localctx = NewWith_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SQLParserRULE_with_clause)

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
			p.Duration_lit()
		}

	case SQLParserT_SHARD:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(203)
			p.Match(SQLParserT_SHARD)
		}
		{
			p.SetState(204)
			p.Shard_num()
		}

	case SQLParserT_REPLICATION:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(205)
			p.Match(SQLParserT_REPLICATION)
		}
		{
			p.SetState(206)
			p.Replica_factor()
		}

	case SQLParserT_TTL:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(207)
			p.Match(SQLParserT_TTL)
		}
		{
			p.SetState(208)
			p.Ttl_val()
		}

	case SQLParserT_META_TTL:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(209)
			p.Match(SQLParserT_META_TTL)
		}
		{
			p.SetState(210)
			p.Metattl_val()
		}

	case SQLParserT_PAST_TTL:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(211)
			p.Match(SQLParserT_PAST_TTL)
		}
		{
			p.SetState(212)
			p.Past_val()
		}

	case SQLParserT_FUTURE_TTL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(213)
			p.Match(SQLParserT_FUTURE_TTL)
		}
		{
			p.SetState(214)
			p.Future_val()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IInterval_define_listContext is an interface to support dynamic dispatch.
type IInterval_define_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInterval_define_listContext differentiates from other interfaces.
	IsInterval_define_listContext()
}

type Interval_define_listContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInterval_define_listContext() *Interval_define_listContext {
	var p = new(Interval_define_listContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_interval_define_list
	return p
}

func (*Interval_define_listContext) IsInterval_define_listContext() {}

func NewInterval_define_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Interval_define_listContext {
	var p = new(Interval_define_listContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_interval_define_list

	return p
}

func (s *Interval_define_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Interval_define_listContext) AllInterval_define() []IInterval_defineContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IInterval_defineContext)(nil)).Elem())
	var tst = make([]IInterval_defineContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IInterval_defineContext)
		}
	}

	return tst
}

func (s *Interval_define_listContext) Interval_define(i int) IInterval_defineContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInterval_defineContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IInterval_defineContext)
}

func (s *Interval_define_listContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *Interval_define_listContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *Interval_define_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Interval_define_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Interval_define_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterInterval_define_list(s)
	}
}

func (s *Interval_define_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitInterval_define_list(s)
	}
}

func (s *Interval_define_listContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitInterval_define_list(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Interval_define_list() (localctx IInterval_define_listContext) {
	localctx = NewInterval_define_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SQLParserRULE_interval_define_list)
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
		p.Interval_define()
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
			p.Interval_define()
		}

		p.SetState(224)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IInterval_defineContext is an interface to support dynamic dispatch.
type IInterval_defineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInterval_defineContext differentiates from other interfaces.
	IsInterval_defineContext()
}

type Interval_defineContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInterval_defineContext() *Interval_defineContext {
	var p = new(Interval_defineContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_interval_define
	return p
}

func (*Interval_defineContext) IsInterval_defineContext() {}

func NewInterval_defineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Interval_defineContext {
	var p = new(Interval_defineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_interval_define

	return p
}

func (s *Interval_defineContext) GetParser() antlr.Parser { return s.parser }

func (s *Interval_defineContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *Interval_defineContext) T_INTERVAL_NAME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL_NAME, 0)
}

func (s *Interval_defineContext) Interval_name_val() IInterval_name_valContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInterval_name_valContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInterval_name_valContext)
}

func (s *Interval_defineContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *Interval_defineContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *Interval_defineContext) T_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TTL, 0)
}

func (s *Interval_defineContext) Ttl_val() ITtl_valContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITtl_valContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITtl_valContext)
}

func (s *Interval_defineContext) T_INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL, 0)
}

func (s *Interval_defineContext) Duration_lit() IDuration_litContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDuration_litContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDuration_litContext)
}

func (s *Interval_defineContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *Interval_defineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Interval_defineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Interval_defineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterInterval_define(s)
	}
}

func (s *Interval_defineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitInterval_define(s)
	}
}

func (s *Interval_defineContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitInterval_define(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Interval_define() (localctx IInterval_defineContext) {
	localctx = NewInterval_defineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SQLParserRULE_interval_define)

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
		p.Interval_name_val()
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
		p.Ttl_val()
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
		p.Duration_lit()
	}
	{
		p.SetState(234)
		p.Match(SQLParserT_CLOSE_P)
	}

	return localctx
}

// IShard_numContext is an interface to support dynamic dispatch.
type IShard_numContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShard_numContext differentiates from other interfaces.
	IsShard_numContext()
}

type Shard_numContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShard_numContext() *Shard_numContext {
	var p = new(Shard_numContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_shard_num
	return p
}

func (*Shard_numContext) IsShard_numContext() {}

func NewShard_numContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Shard_numContext {
	var p = new(Shard_numContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_shard_num

	return p
}

func (s *Shard_numContext) GetParser() antlr.Parser { return s.parser }

func (s *Shard_numContext) Int_number() IInt_numberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInt_numberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInt_numberContext)
}

func (s *Shard_numContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Shard_numContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Shard_numContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShard_num(s)
	}
}

func (s *Shard_numContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShard_num(s)
	}
}

func (s *Shard_numContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShard_num(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Shard_num() (localctx IShard_numContext) {
	localctx = NewShard_numContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SQLParserRULE_shard_num)

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
		p.Int_number()
	}

	return localctx
}

// ITtl_valContext is an interface to support dynamic dispatch.
type ITtl_valContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTtl_valContext differentiates from other interfaces.
	IsTtl_valContext()
}

type Ttl_valContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTtl_valContext() *Ttl_valContext {
	var p = new(Ttl_valContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_ttl_val
	return p
}

func (*Ttl_valContext) IsTtl_valContext() {}

func NewTtl_valContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Ttl_valContext {
	var p = new(Ttl_valContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_ttl_val

	return p
}

func (s *Ttl_valContext) GetParser() antlr.Parser { return s.parser }

func (s *Ttl_valContext) Duration_lit() IDuration_litContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDuration_litContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDuration_litContext)
}

func (s *Ttl_valContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Ttl_valContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Ttl_valContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTtl_val(s)
	}
}

func (s *Ttl_valContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTtl_val(s)
	}
}

func (s *Ttl_valContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTtl_val(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Ttl_val() (localctx ITtl_valContext) {
	localctx = NewTtl_valContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SQLParserRULE_ttl_val)

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
		p.Duration_lit()
	}

	return localctx
}

// IMetattl_valContext is an interface to support dynamic dispatch.
type IMetattl_valContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMetattl_valContext differentiates from other interfaces.
	IsMetattl_valContext()
}

type Metattl_valContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetattl_valContext() *Metattl_valContext {
	var p = new(Metattl_valContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_metattl_val
	return p
}

func (*Metattl_valContext) IsMetattl_valContext() {}

func NewMetattl_valContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Metattl_valContext {
	var p = new(Metattl_valContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_metattl_val

	return p
}

func (s *Metattl_valContext) GetParser() antlr.Parser { return s.parser }

func (s *Metattl_valContext) Duration_lit() IDuration_litContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDuration_litContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDuration_litContext)
}

func (s *Metattl_valContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Metattl_valContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Metattl_valContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterMetattl_val(s)
	}
}

func (s *Metattl_valContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitMetattl_val(s)
	}
}

func (s *Metattl_valContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitMetattl_val(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Metattl_val() (localctx IMetattl_valContext) {
	localctx = NewMetattl_valContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SQLParserRULE_metattl_val)

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
		p.Duration_lit()
	}

	return localctx
}

// IPast_valContext is an interface to support dynamic dispatch.
type IPast_valContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPast_valContext differentiates from other interfaces.
	IsPast_valContext()
}

type Past_valContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPast_valContext() *Past_valContext {
	var p = new(Past_valContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_past_val
	return p
}

func (*Past_valContext) IsPast_valContext() {}

func NewPast_valContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Past_valContext {
	var p = new(Past_valContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_past_val

	return p
}

func (s *Past_valContext) GetParser() antlr.Parser { return s.parser }

func (s *Past_valContext) Duration_lit() IDuration_litContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDuration_litContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDuration_litContext)
}

func (s *Past_valContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Past_valContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Past_valContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterPast_val(s)
	}
}

func (s *Past_valContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitPast_val(s)
	}
}

func (s *Past_valContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitPast_val(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Past_val() (localctx IPast_valContext) {
	localctx = NewPast_valContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SQLParserRULE_past_val)

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
		p.Duration_lit()
	}

	return localctx
}

// IFuture_valContext is an interface to support dynamic dispatch.
type IFuture_valContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFuture_valContext differentiates from other interfaces.
	IsFuture_valContext()
}

type Future_valContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuture_valContext() *Future_valContext {
	var p = new(Future_valContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_future_val
	return p
}

func (*Future_valContext) IsFuture_valContext() {}

func NewFuture_valContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Future_valContext {
	var p = new(Future_valContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_future_val

	return p
}

func (s *Future_valContext) GetParser() antlr.Parser { return s.parser }

func (s *Future_valContext) Duration_lit() IDuration_litContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDuration_litContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDuration_litContext)
}

func (s *Future_valContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Future_valContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Future_valContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFuture_val(s)
	}
}

func (s *Future_valContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFuture_val(s)
	}
}

func (s *Future_valContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFuture_val(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Future_val() (localctx IFuture_valContext) {
	localctx = NewFuture_valContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SQLParserRULE_future_val)

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
		p.Duration_lit()
	}

	return localctx
}

// IInterval_name_valContext is an interface to support dynamic dispatch.
type IInterval_name_valContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInterval_name_valContext differentiates from other interfaces.
	IsInterval_name_valContext()
}

type Interval_name_valContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInterval_name_valContext() *Interval_name_valContext {
	var p = new(Interval_name_valContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_interval_name_val
	return p
}

func (*Interval_name_valContext) IsInterval_name_valContext() {}

func NewInterval_name_valContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Interval_name_valContext {
	var p = new(Interval_name_valContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_interval_name_val

	return p
}

func (s *Interval_name_valContext) GetParser() antlr.Parser { return s.parser }

func (s *Interval_name_valContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *Interval_name_valContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Interval_name_valContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Interval_name_valContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterInterval_name_val(s)
	}
}

func (s *Interval_name_valContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitInterval_name_val(s)
	}
}

func (s *Interval_name_valContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitInterval_name_val(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Interval_name_val() (localctx IInterval_name_valContext) {
	localctx = NewInterval_name_valContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, SQLParserRULE_interval_name_val)

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

// IReplica_factorContext is an interface to support dynamic dispatch.
type IReplica_factorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsReplica_factorContext differentiates from other interfaces.
	IsReplica_factorContext()
}

type Replica_factorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReplica_factorContext() *Replica_factorContext {
	var p = new(Replica_factorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_replica_factor
	return p
}

func (*Replica_factorContext) IsReplica_factorContext() {}

func NewReplica_factorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Replica_factorContext {
	var p = new(Replica_factorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_replica_factor

	return p
}

func (s *Replica_factorContext) GetParser() antlr.Parser { return s.parser }

func (s *Replica_factorContext) Int_number() IInt_numberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInt_numberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInt_numberContext)
}

func (s *Replica_factorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Replica_factorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Replica_factorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterReplica_factor(s)
	}
}

func (s *Replica_factorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitReplica_factor(s)
	}
}

func (s *Replica_factorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitReplica_factor(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Replica_factor() (localctx IReplica_factorContext) {
	localctx = NewReplica_factorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, SQLParserRULE_replica_factor)

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
		p.Int_number()
	}

	return localctx
}

// IDatabase_nameContext is an interface to support dynamic dispatch.
type IDatabase_nameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDatabase_nameContext differentiates from other interfaces.
	IsDatabase_nameContext()
}

type Database_nameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDatabase_nameContext() *Database_nameContext {
	var p = new(Database_nameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_database_name
	return p
}

func (*Database_nameContext) IsDatabase_nameContext() {}

func NewDatabase_nameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Database_nameContext {
	var p = new(Database_nameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_database_name

	return p
}

func (s *Database_nameContext) GetParser() antlr.Parser { return s.parser }

func (s *Database_nameContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *Database_nameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Database_nameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Database_nameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDatabase_name(s)
	}
}

func (s *Database_nameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDatabase_name(s)
	}
}

func (s *Database_nameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDatabase_name(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Database_name() (localctx IDatabase_nameContext) {
	localctx = NewDatabase_nameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, SQLParserRULE_database_name)

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

// IUpdate_database_stmtContext is an interface to support dynamic dispatch.
type IUpdate_database_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUpdate_database_stmtContext differentiates from other interfaces.
	IsUpdate_database_stmtContext()
}

type Update_database_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUpdate_database_stmtContext() *Update_database_stmtContext {
	var p = new(Update_database_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_update_database_stmt
	return p
}

func (*Update_database_stmtContext) IsUpdate_database_stmtContext() {}

func NewUpdate_database_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Update_database_stmtContext {
	var p = new(Update_database_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_update_database_stmt

	return p
}

func (s *Update_database_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Update_database_stmtContext) T_UPDATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_UPDATE, 0)
}

func (s *Update_database_stmtContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *Update_database_stmtContext) Database_name() IDatabase_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDatabase_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDatabase_nameContext)
}

func (s *Update_database_stmtContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *Update_database_stmtContext) With_clause_list() IWith_clause_listContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWith_clause_listContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWith_clause_listContext)
}

func (s *Update_database_stmtContext) T_COMMA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, 0)
}

func (s *Update_database_stmtContext) Interval_define_list() IInterval_define_listContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInterval_define_listContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInterval_define_listContext)
}

func (s *Update_database_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Update_database_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Update_database_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterUpdate_database_stmt(s)
	}
}

func (s *Update_database_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitUpdate_database_stmt(s)
	}
}

func (s *Update_database_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitUpdate_database_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Update_database_stmt() (localctx IUpdate_database_stmtContext) {
	localctx = NewUpdate_database_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, SQLParserRULE_update_database_stmt)
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
		p.Database_name()
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
			p.With_clause_list()
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
			p.Interval_define_list()
		}

	}

	return localctx
}

// IDrop_database_stmtContext is an interface to support dynamic dispatch.
type IDrop_database_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDrop_database_stmtContext differentiates from other interfaces.
	IsDrop_database_stmtContext()
}

type Drop_database_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDrop_database_stmtContext() *Drop_database_stmtContext {
	var p = new(Drop_database_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_drop_database_stmt
	return p
}

func (*Drop_database_stmtContext) IsDrop_database_stmtContext() {}

func NewDrop_database_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Drop_database_stmtContext {
	var p = new(Drop_database_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_drop_database_stmt

	return p
}

func (s *Drop_database_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Drop_database_stmtContext) T_DROP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DROP, 0)
}

func (s *Drop_database_stmtContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *Drop_database_stmtContext) Database_name() IDatabase_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDatabase_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDatabase_nameContext)
}

func (s *Drop_database_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Drop_database_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Drop_database_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDrop_database_stmt(s)
	}
}

func (s *Drop_database_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDrop_database_stmt(s)
	}
}

func (s *Drop_database_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDrop_database_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Drop_database_stmt() (localctx IDrop_database_stmtContext) {
	localctx = NewDrop_database_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, SQLParserRULE_drop_database_stmt)

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
		p.Database_name()
	}

	return localctx
}

// IShow_databases_stmtContext is an interface to support dynamic dispatch.
type IShow_databases_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShow_databases_stmtContext differentiates from other interfaces.
	IsShow_databases_stmtContext()
}

type Show_databases_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShow_databases_stmtContext() *Show_databases_stmtContext {
	var p = new(Show_databases_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_show_databases_stmt
	return p
}

func (*Show_databases_stmtContext) IsShow_databases_stmtContext() {}

func NewShow_databases_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Show_databases_stmtContext {
	var p = new(Show_databases_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_show_databases_stmt

	return p
}

func (s *Show_databases_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Show_databases_stmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Show_databases_stmtContext) T_DATASBAES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAES, 0)
}

func (s *Show_databases_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Show_databases_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Show_databases_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShow_databases_stmt(s)
	}
}

func (s *Show_databases_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShow_databases_stmt(s)
	}
}

func (s *Show_databases_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShow_databases_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Show_databases_stmt() (localctx IShow_databases_stmtContext) {
	localctx = NewShow_databases_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, SQLParserRULE_show_databases_stmt)

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

// IShow_node_stmtContext is an interface to support dynamic dispatch.
type IShow_node_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShow_node_stmtContext differentiates from other interfaces.
	IsShow_node_stmtContext()
}

type Show_node_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShow_node_stmtContext() *Show_node_stmtContext {
	var p = new(Show_node_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_show_node_stmt
	return p
}

func (*Show_node_stmtContext) IsShow_node_stmtContext() {}

func NewShow_node_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Show_node_stmtContext {
	var p = new(Show_node_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_show_node_stmt

	return p
}

func (s *Show_node_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Show_node_stmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Show_node_stmtContext) T_NODE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NODE, 0)
}

func (s *Show_node_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Show_node_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Show_node_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShow_node_stmt(s)
	}
}

func (s *Show_node_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShow_node_stmt(s)
	}
}

func (s *Show_node_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShow_node_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Show_node_stmt() (localctx IShow_node_stmtContext) {
	localctx = NewShow_node_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, SQLParserRULE_show_node_stmt)

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

// IShow_measurements_stmtContext is an interface to support dynamic dispatch.
type IShow_measurements_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShow_measurements_stmtContext differentiates from other interfaces.
	IsShow_measurements_stmtContext()
}

type Show_measurements_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShow_measurements_stmtContext() *Show_measurements_stmtContext {
	var p = new(Show_measurements_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_show_measurements_stmt
	return p
}

func (*Show_measurements_stmtContext) IsShow_measurements_stmtContext() {}

func NewShow_measurements_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Show_measurements_stmtContext {
	var p = new(Show_measurements_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_show_measurements_stmt

	return p
}

func (s *Show_measurements_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Show_measurements_stmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Show_measurements_stmtContext) T_MEASUREMENTS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MEASUREMENTS, 0)
}

func (s *Show_measurements_stmtContext) With_measurement_clause() IWith_measurement_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWith_measurement_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWith_measurement_clauseContext)
}

func (s *Show_measurements_stmtContext) Limit_clause() ILimit_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimit_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimit_clauseContext)
}

func (s *Show_measurements_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Show_measurements_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Show_measurements_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShow_measurements_stmt(s)
	}
}

func (s *Show_measurements_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShow_measurements_stmt(s)
	}
}

func (s *Show_measurements_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShow_measurements_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Show_measurements_stmt() (localctx IShow_measurements_stmtContext) {
	localctx = NewShow_measurements_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, SQLParserRULE_show_measurements_stmt)
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
			p.With_measurement_clause()
		}

	}
	p.SetState(279)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(278)
			p.Limit_clause()
		}

	}

	return localctx
}

// IShow_tag_keys_stmtContext is an interface to support dynamic dispatch.
type IShow_tag_keys_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShow_tag_keys_stmtContext differentiates from other interfaces.
	IsShow_tag_keys_stmtContext()
}

type Show_tag_keys_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShow_tag_keys_stmtContext() *Show_tag_keys_stmtContext {
	var p = new(Show_tag_keys_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_show_tag_keys_stmt
	return p
}

func (*Show_tag_keys_stmtContext) IsShow_tag_keys_stmtContext() {}

func NewShow_tag_keys_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Show_tag_keys_stmtContext {
	var p = new(Show_tag_keys_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_show_tag_keys_stmt

	return p
}

func (s *Show_tag_keys_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Show_tag_keys_stmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Show_tag_keys_stmtContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *Show_tag_keys_stmtContext) T_KEYS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEYS, 0)
}

func (s *Show_tag_keys_stmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *Show_tag_keys_stmtContext) Metric_name() IMetric_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetric_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetric_nameContext)
}

func (s *Show_tag_keys_stmtContext) Limit_clause() ILimit_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimit_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimit_clauseContext)
}

func (s *Show_tag_keys_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Show_tag_keys_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Show_tag_keys_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShow_tag_keys_stmt(s)
	}
}

func (s *Show_tag_keys_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShow_tag_keys_stmt(s)
	}
}

func (s *Show_tag_keys_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShow_tag_keys_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Show_tag_keys_stmt() (localctx IShow_tag_keys_stmtContext) {
	localctx = NewShow_tag_keys_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, SQLParserRULE_show_tag_keys_stmt)
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
		p.Metric_name()
	}
	p.SetState(287)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(286)
			p.Limit_clause()
		}

	}

	return localctx
}

// IShow_info_stmtContext is an interface to support dynamic dispatch.
type IShow_info_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShow_info_stmtContext differentiates from other interfaces.
	IsShow_info_stmtContext()
}

type Show_info_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShow_info_stmtContext() *Show_info_stmtContext {
	var p = new(Show_info_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_show_info_stmt
	return p
}

func (*Show_info_stmtContext) IsShow_info_stmtContext() {}

func NewShow_info_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Show_info_stmtContext {
	var p = new(Show_info_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_show_info_stmt

	return p
}

func (s *Show_info_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Show_info_stmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Show_info_stmtContext) T_INFO() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INFO, 0)
}

func (s *Show_info_stmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *Show_info_stmtContext) Metric_name() IMetric_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetric_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetric_nameContext)
}

func (s *Show_info_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Show_info_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Show_info_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShow_info_stmt(s)
	}
}

func (s *Show_info_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShow_info_stmt(s)
	}
}

func (s *Show_info_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShow_info_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Show_info_stmt() (localctx IShow_info_stmtContext) {
	localctx = NewShow_info_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, SQLParserRULE_show_info_stmt)

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
		p.Metric_name()
	}

	return localctx
}

// IShow_tag_values_stmtContext is an interface to support dynamic dispatch.
type IShow_tag_values_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShow_tag_values_stmtContext differentiates from other interfaces.
	IsShow_tag_values_stmtContext()
}

type Show_tag_values_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShow_tag_values_stmtContext() *Show_tag_values_stmtContext {
	var p = new(Show_tag_values_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_show_tag_values_stmt
	return p
}

func (*Show_tag_values_stmtContext) IsShow_tag_values_stmtContext() {}

func NewShow_tag_values_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Show_tag_values_stmtContext {
	var p = new(Show_tag_values_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_show_tag_values_stmt

	return p
}

func (s *Show_tag_values_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Show_tag_values_stmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Show_tag_values_stmtContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *Show_tag_values_stmtContext) T_VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUES, 0)
}

func (s *Show_tag_values_stmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *Show_tag_values_stmtContext) Metric_name() IMetric_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetric_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetric_nameContext)
}

func (s *Show_tag_values_stmtContext) With_tag_clause() IWith_tag_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWith_tag_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWith_tag_clauseContext)
}

func (s *Show_tag_values_stmtContext) Where_tag_cascade() IWhere_tag_cascadeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWhere_tag_cascadeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWhere_tag_cascadeContext)
}

func (s *Show_tag_values_stmtContext) Limit_clause() ILimit_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimit_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimit_clauseContext)
}

func (s *Show_tag_values_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Show_tag_values_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Show_tag_values_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShow_tag_values_stmt(s)
	}
}

func (s *Show_tag_values_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShow_tag_values_stmt(s)
	}
}

func (s *Show_tag_values_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShow_tag_values_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Show_tag_values_stmt() (localctx IShow_tag_values_stmtContext) {
	localctx = NewShow_tag_values_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, SQLParserRULE_show_tag_values_stmt)
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
		p.Metric_name()
	}
	{
		p.SetState(299)
		p.With_tag_clause()
	}
	p.SetState(301)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_WHERE {
		{
			p.SetState(300)
			p.Where_tag_cascade()
		}

	}
	p.SetState(304)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(303)
			p.Limit_clause()
		}

	}

	return localctx
}

// IShow_tag_values_info_stmtContext is an interface to support dynamic dispatch.
type IShow_tag_values_info_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShow_tag_values_info_stmtContext differentiates from other interfaces.
	IsShow_tag_values_info_stmtContext()
}

type Show_tag_values_info_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShow_tag_values_info_stmtContext() *Show_tag_values_info_stmtContext {
	var p = new(Show_tag_values_info_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_show_tag_values_info_stmt
	return p
}

func (*Show_tag_values_info_stmtContext) IsShow_tag_values_info_stmtContext() {}

func NewShow_tag_values_info_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Show_tag_values_info_stmtContext {
	var p = new(Show_tag_values_info_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_show_tag_values_info_stmt

	return p
}

func (s *Show_tag_values_info_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Show_tag_values_info_stmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Show_tag_values_info_stmtContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *Show_tag_values_info_stmtContext) T_VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUES, 0)
}

func (s *Show_tag_values_info_stmtContext) T_INFO() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INFO, 0)
}

func (s *Show_tag_values_info_stmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *Show_tag_values_info_stmtContext) Metric_name() IMetric_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetric_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetric_nameContext)
}

func (s *Show_tag_values_info_stmtContext) With_tag_clause() IWith_tag_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWith_tag_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWith_tag_clauseContext)
}

func (s *Show_tag_values_info_stmtContext) Where_tag_cascade() IWhere_tag_cascadeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWhere_tag_cascadeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWhere_tag_cascadeContext)
}

func (s *Show_tag_values_info_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Show_tag_values_info_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Show_tag_values_info_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShow_tag_values_info_stmt(s)
	}
}

func (s *Show_tag_values_info_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShow_tag_values_info_stmt(s)
	}
}

func (s *Show_tag_values_info_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShow_tag_values_info_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Show_tag_values_info_stmt() (localctx IShow_tag_values_info_stmtContext) {
	localctx = NewShow_tag_values_info_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, SQLParserRULE_show_tag_values_info_stmt)

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
		p.Metric_name()
	}
	{
		p.SetState(312)
		p.With_tag_clause()
	}
	{
		p.SetState(313)
		p.Where_tag_cascade()
	}

	return localctx
}

// IShow_field_keys_stmtContext is an interface to support dynamic dispatch.
type IShow_field_keys_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShow_field_keys_stmtContext differentiates from other interfaces.
	IsShow_field_keys_stmtContext()
}

type Show_field_keys_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShow_field_keys_stmtContext() *Show_field_keys_stmtContext {
	var p = new(Show_field_keys_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_show_field_keys_stmt
	return p
}

func (*Show_field_keys_stmtContext) IsShow_field_keys_stmtContext() {}

func NewShow_field_keys_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Show_field_keys_stmtContext {
	var p = new(Show_field_keys_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_show_field_keys_stmt

	return p
}

func (s *Show_field_keys_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Show_field_keys_stmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Show_field_keys_stmtContext) T_FIELD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIELD, 0)
}

func (s *Show_field_keys_stmtContext) T_KEYS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEYS, 0)
}

func (s *Show_field_keys_stmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *Show_field_keys_stmtContext) Metric_name() IMetric_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetric_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetric_nameContext)
}

func (s *Show_field_keys_stmtContext) Limit_clause() ILimit_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimit_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimit_clauseContext)
}

func (s *Show_field_keys_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Show_field_keys_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Show_field_keys_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShow_field_keys_stmt(s)
	}
}

func (s *Show_field_keys_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShow_field_keys_stmt(s)
	}
}

func (s *Show_field_keys_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShow_field_keys_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Show_field_keys_stmt() (localctx IShow_field_keys_stmtContext) {
	localctx = NewShow_field_keys_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, SQLParserRULE_show_field_keys_stmt)
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
		p.Metric_name()
	}
	p.SetState(321)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(320)
			p.Limit_clause()
		}

	}

	return localctx
}

// IShow_queries_stmtContext is an interface to support dynamic dispatch.
type IShow_queries_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShow_queries_stmtContext differentiates from other interfaces.
	IsShow_queries_stmtContext()
}

type Show_queries_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShow_queries_stmtContext() *Show_queries_stmtContext {
	var p = new(Show_queries_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_show_queries_stmt
	return p
}

func (*Show_queries_stmtContext) IsShow_queries_stmtContext() {}

func NewShow_queries_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Show_queries_stmtContext {
	var p = new(Show_queries_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_show_queries_stmt

	return p
}

func (s *Show_queries_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Show_queries_stmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Show_queries_stmtContext) T_QUERIES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUERIES, 0)
}

func (s *Show_queries_stmtContext) Limit_clause() ILimit_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimit_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimit_clauseContext)
}

func (s *Show_queries_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Show_queries_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Show_queries_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShow_queries_stmt(s)
	}
}

func (s *Show_queries_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShow_queries_stmt(s)
	}
}

func (s *Show_queries_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShow_queries_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Show_queries_stmt() (localctx IShow_queries_stmtContext) {
	localctx = NewShow_queries_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, SQLParserRULE_show_queries_stmt)
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
			p.Limit_clause()
		}

	}

	return localctx
}

// IShow_stats_stmtContext is an interface to support dynamic dispatch.
type IShow_stats_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShow_stats_stmtContext differentiates from other interfaces.
	IsShow_stats_stmtContext()
}

type Show_stats_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShow_stats_stmtContext() *Show_stats_stmtContext {
	var p = new(Show_stats_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_show_stats_stmt
	return p
}

func (*Show_stats_stmtContext) IsShow_stats_stmtContext() {}

func NewShow_stats_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Show_stats_stmtContext {
	var p = new(Show_stats_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_show_stats_stmt

	return p
}

func (s *Show_stats_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Show_stats_stmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Show_stats_stmtContext) T_STATS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATS, 0)
}

func (s *Show_stats_stmtContext) T_FOR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FOR, 0)
}

func (s *Show_stats_stmtContext) Module() IModuleContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IModuleContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IModuleContext)
}

func (s *Show_stats_stmtContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *Show_stats_stmtContext) Component() IComponentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IComponentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IComponentContext)
}

func (s *Show_stats_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Show_stats_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Show_stats_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShow_stats_stmt(s)
	}
}

func (s *Show_stats_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShow_stats_stmt(s)
	}
}

func (s *Show_stats_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShow_stats_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Show_stats_stmt() (localctx IShow_stats_stmtContext) {
	localctx = NewShow_stats_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, SQLParserRULE_show_stats_stmt)
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

// IWith_measurement_clauseContext is an interface to support dynamic dispatch.
type IWith_measurement_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWith_measurement_clauseContext differentiates from other interfaces.
	IsWith_measurement_clauseContext()
}

type With_measurement_clauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWith_measurement_clauseContext() *With_measurement_clauseContext {
	var p = new(With_measurement_clauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_with_measurement_clause
	return p
}

func (*With_measurement_clauseContext) IsWith_measurement_clauseContext() {}

func NewWith_measurement_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *With_measurement_clauseContext {
	var p = new(With_measurement_clauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_with_measurement_clause

	return p
}

func (s *With_measurement_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *With_measurement_clauseContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *With_measurement_clauseContext) T_MEASUREMENT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MEASUREMENT, 0)
}

func (s *With_measurement_clauseContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *With_measurement_clauseContext) Metric_name() IMetric_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetric_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetric_nameContext)
}

func (s *With_measurement_clauseContext) T_REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REGEXP, 0)
}

func (s *With_measurement_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *With_measurement_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *With_measurement_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWith_measurement_clause(s)
	}
}

func (s *With_measurement_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWith_measurement_clause(s)
	}
}

func (s *With_measurement_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWith_measurement_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) With_measurement_clause() (localctx IWith_measurement_clauseContext) {
	localctx = NewWith_measurement_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, SQLParserRULE_with_measurement_clause)

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
			p.Metric_name()
		}

	case SQLParserT_REGEXP:
		{
			p.SetState(342)
			p.Match(SQLParserT_REGEXP)
		}
		{
			p.SetState(343)
			p.Metric_name()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IWith_tag_clauseContext is an interface to support dynamic dispatch.
type IWith_tag_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWith_tag_clauseContext differentiates from other interfaces.
	IsWith_tag_clauseContext()
}

type With_tag_clauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWith_tag_clauseContext() *With_tag_clauseContext {
	var p = new(With_tag_clauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_with_tag_clause
	return p
}

func (*With_tag_clauseContext) IsWith_tag_clauseContext() {}

func NewWith_tag_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *With_tag_clauseContext {
	var p = new(With_tag_clauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_with_tag_clause

	return p
}

func (s *With_tag_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *With_tag_clauseContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *With_tag_clauseContext) T_KEY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEY, 0)
}

func (s *With_tag_clauseContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *With_tag_clauseContext) Tag_key() ITag_keyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_keyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_keyContext)
}

func (s *With_tag_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *With_tag_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *With_tag_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWith_tag_clause(s)
	}
}

func (s *With_tag_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWith_tag_clause(s)
	}
}

func (s *With_tag_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWith_tag_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) With_tag_clause() (localctx IWith_tag_clauseContext) {
	localctx = NewWith_tag_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, SQLParserRULE_with_tag_clause)

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
		p.Tag_key()
	}

	return localctx
}

// IWhere_tag_cascadeContext is an interface to support dynamic dispatch.
type IWhere_tag_cascadeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWhere_tag_cascadeContext differentiates from other interfaces.
	IsWhere_tag_cascadeContext()
}

type Where_tag_cascadeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhere_tag_cascadeContext() *Where_tag_cascadeContext {
	var p = new(Where_tag_cascadeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_where_tag_cascade
	return p
}

func (*Where_tag_cascadeContext) IsWhere_tag_cascadeContext() {}

func NewWhere_tag_cascadeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Where_tag_cascadeContext {
	var p = new(Where_tag_cascadeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_where_tag_cascade

	return p
}

func (s *Where_tag_cascadeContext) GetParser() antlr.Parser { return s.parser }

func (s *Where_tag_cascadeContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *Where_tag_cascadeContext) Tag_cascade_expr() ITag_cascade_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_cascade_exprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_cascade_exprContext)
}

func (s *Where_tag_cascadeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Where_tag_cascadeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Where_tag_cascadeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWhere_tag_cascade(s)
	}
}

func (s *Where_tag_cascadeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWhere_tag_cascade(s)
	}
}

func (s *Where_tag_cascadeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWhere_tag_cascade(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Where_tag_cascade() (localctx IWhere_tag_cascadeContext) {
	localctx = NewWhere_tag_cascadeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, SQLParserRULE_where_tag_cascade)

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
		p.Tag_cascade_expr()
	}

	return localctx
}

// IKill_query_stmtContext is an interface to support dynamic dispatch.
type IKill_query_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsKill_query_stmtContext differentiates from other interfaces.
	IsKill_query_stmtContext()
}

type Kill_query_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyKill_query_stmtContext() *Kill_query_stmtContext {
	var p = new(Kill_query_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_kill_query_stmt
	return p
}

func (*Kill_query_stmtContext) IsKill_query_stmtContext() {}

func NewKill_query_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Kill_query_stmtContext {
	var p = new(Kill_query_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_kill_query_stmt

	return p
}

func (s *Kill_query_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Kill_query_stmtContext) T_KILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KILL, 0)
}

func (s *Kill_query_stmtContext) T_QUERY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUERY, 0)
}

func (s *Kill_query_stmtContext) Query_id() IQuery_idContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQuery_idContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IQuery_idContext)
}

func (s *Kill_query_stmtContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *Kill_query_stmtContext) Server_id() IServer_idContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IServer_idContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IServer_idContext)
}

func (s *Kill_query_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Kill_query_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Kill_query_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterKill_query_stmt(s)
	}
}

func (s *Kill_query_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitKill_query_stmt(s)
	}
}

func (s *Kill_query_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitKill_query_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Kill_query_stmt() (localctx IKill_query_stmtContext) {
	localctx = NewKill_query_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, SQLParserRULE_kill_query_stmt)
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
		p.Query_id()
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
			p.Server_id()
		}

	}

	return localctx
}

// IQuery_idContext is an interface to support dynamic dispatch.
type IQuery_idContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQuery_idContext differentiates from other interfaces.
	IsQuery_idContext()
}

type Query_idContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQuery_idContext() *Query_idContext {
	var p = new(Query_idContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_query_id
	return p
}

func (*Query_idContext) IsQuery_idContext() {}

func NewQuery_idContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Query_idContext {
	var p = new(Query_idContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_query_id

	return p
}

func (s *Query_idContext) GetParser() antlr.Parser { return s.parser }

func (s *Query_idContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *Query_idContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Query_idContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Query_idContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterQuery_id(s)
	}
}

func (s *Query_idContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitQuery_id(s)
	}
}

func (s *Query_idContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitQuery_id(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Query_id() (localctx IQuery_idContext) {
	localctx = NewQuery_idContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, SQLParserRULE_query_id)

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

// IServer_idContext is an interface to support dynamic dispatch.
type IServer_idContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsServer_idContext differentiates from other interfaces.
	IsServer_idContext()
}

type Server_idContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyServer_idContext() *Server_idContext {
	var p = new(Server_idContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_server_id
	return p
}

func (*Server_idContext) IsServer_idContext() {}

func NewServer_idContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Server_idContext {
	var p = new(Server_idContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_server_id

	return p
}

func (s *Server_idContext) GetParser() antlr.Parser { return s.parser }

func (s *Server_idContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *Server_idContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Server_idContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Server_idContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterServer_id(s)
	}
}

func (s *Server_idContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitServer_id(s)
	}
}

func (s *Server_idContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitServer_id(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Server_id() (localctx IServer_idContext) {
	localctx = NewServer_idContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, SQLParserRULE_server_id)

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

// IQuery_stmtContext is an interface to support dynamic dispatch.
type IQuery_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQuery_stmtContext differentiates from other interfaces.
	IsQuery_stmtContext()
}

type Query_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQuery_stmtContext() *Query_stmtContext {
	var p = new(Query_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_query_stmt
	return p
}

func (*Query_stmtContext) IsQuery_stmtContext() {}

func NewQuery_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Query_stmtContext {
	var p = new(Query_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_query_stmt

	return p
}

func (s *Query_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Query_stmtContext) T_SELECT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SELECT, 0)
}

func (s *Query_stmtContext) Fields() IFieldsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldsContext)
}

func (s *Query_stmtContext) From_clause() IFrom_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFrom_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFrom_clauseContext)
}

func (s *Query_stmtContext) T_EXPLAIN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EXPLAIN, 0)
}

func (s *Query_stmtContext) Where_clause() IWhere_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWhere_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWhere_clauseContext)
}

func (s *Query_stmtContext) Group_by_clause() IGroup_by_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IGroup_by_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IGroup_by_clauseContext)
}

func (s *Query_stmtContext) Interval_by_clause() IInterval_by_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInterval_by_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInterval_by_clauseContext)
}

func (s *Query_stmtContext) Order_by_clause() IOrder_by_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOrder_by_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOrder_by_clauseContext)
}

func (s *Query_stmtContext) Limit_clause() ILimit_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimit_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimit_clauseContext)
}

func (s *Query_stmtContext) T_WITH_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH_VALUE, 0)
}

func (s *Query_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Query_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Query_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterQuery_stmt(s)
	}
}

func (s *Query_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitQuery_stmt(s)
	}
}

func (s *Query_stmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitQuery_stmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Query_stmt() (localctx IQuery_stmtContext) {
	localctx = NewQuery_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, SQLParserRULE_query_stmt)
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
		p.From_clause()
	}
	p.SetState(376)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_WHERE {
		{
			p.SetState(375)
			p.Where_clause()
		}

	}
	p.SetState(379)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_GROUP {
		{
			p.SetState(378)
			p.Group_by_clause()
		}

	}
	p.SetState(382)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_INTERVAL {
		{
			p.SetState(381)
			p.Interval_by_clause()
		}

	}
	p.SetState(385)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_ORDER {
		{
			p.SetState(384)
			p.Order_by_clause()
		}

	}
	p.SetState(388)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserT_LIMIT {
		{
			p.SetState(387)
			p.Limit_clause()
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

// IFrom_clauseContext is an interface to support dynamic dispatch.
type IFrom_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFrom_clauseContext differentiates from other interfaces.
	IsFrom_clauseContext()
}

type From_clauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFrom_clauseContext() *From_clauseContext {
	var p = new(From_clauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_from_clause
	return p
}

func (*From_clauseContext) IsFrom_clauseContext() {}

func NewFrom_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *From_clauseContext {
	var p = new(From_clauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_from_clause

	return p
}

func (s *From_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *From_clauseContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *From_clauseContext) Metric_name() IMetric_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetric_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetric_nameContext)
}

func (s *From_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *From_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *From_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFrom_clause(s)
	}
}

func (s *From_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFrom_clause(s)
	}
}

func (s *From_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFrom_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) From_clause() (localctx IFrom_clauseContext) {
	localctx = NewFrom_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, SQLParserRULE_from_clause)

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
		p.Metric_name()
	}

	return localctx
}

// IWhere_clauseContext is an interface to support dynamic dispatch.
type IWhere_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWhere_clauseContext differentiates from other interfaces.
	IsWhere_clauseContext()
}

type Where_clauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhere_clauseContext() *Where_clauseContext {
	var p = new(Where_clauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_where_clause
	return p
}

func (*Where_clauseContext) IsWhere_clauseContext() {}

func NewWhere_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Where_clauseContext {
	var p = new(Where_clauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_where_clause

	return p
}

func (s *Where_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *Where_clauseContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *Where_clauseContext) Clause_boolean_expr() IClause_boolean_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IClause_boolean_exprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IClause_boolean_exprContext)
}

func (s *Where_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Where_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Where_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWhere_clause(s)
	}
}

func (s *Where_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWhere_clause(s)
	}
}

func (s *Where_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWhere_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Where_clause() (localctx IWhere_clauseContext) {
	localctx = NewWhere_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, SQLParserRULE_where_clause)

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
		p.clause_boolean_expr(0)
	}

	return localctx
}

// IClause_boolean_exprContext is an interface to support dynamic dispatch.
type IClause_boolean_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsClause_boolean_exprContext differentiates from other interfaces.
	IsClause_boolean_exprContext()
}

type Clause_boolean_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyClause_boolean_exprContext() *Clause_boolean_exprContext {
	var p = new(Clause_boolean_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_clause_boolean_expr
	return p
}

func (*Clause_boolean_exprContext) IsClause_boolean_exprContext() {}

func NewClause_boolean_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Clause_boolean_exprContext {
	var p = new(Clause_boolean_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_clause_boolean_expr

	return p
}

func (s *Clause_boolean_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Clause_boolean_exprContext) Tag_boolean_expr() ITag_boolean_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_boolean_exprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_boolean_exprContext)
}

func (s *Clause_boolean_exprContext) Time_expr() ITime_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITime_exprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITime_exprContext)
}

func (s *Clause_boolean_exprContext) AllClause_boolean_expr() []IClause_boolean_exprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IClause_boolean_exprContext)(nil)).Elem())
	var tst = make([]IClause_boolean_exprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IClause_boolean_exprContext)
		}
	}

	return tst
}

func (s *Clause_boolean_exprContext) Clause_boolean_expr(i int) IClause_boolean_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IClause_boolean_exprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IClause_boolean_exprContext)
}

func (s *Clause_boolean_exprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *Clause_boolean_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Clause_boolean_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Clause_boolean_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterClause_boolean_expr(s)
	}
}

func (s *Clause_boolean_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitClause_boolean_expr(s)
	}
}

func (s *Clause_boolean_exprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitClause_boolean_expr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Clause_boolean_expr() (localctx IClause_boolean_exprContext) {
	return p.clause_boolean_expr(0)
}

func (p *SQLParser) clause_boolean_expr(_p int) (localctx IClause_boolean_exprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewClause_boolean_exprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IClause_boolean_exprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 82
	p.EnterRecursionRule(localctx, 82, SQLParserRULE_clause_boolean_expr, _p)

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
			p.tag_boolean_expr(0)
		}

	case 2:
		{
			p.SetState(416)
			p.Time_expr()
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
			localctx = NewClause_boolean_exprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_clause_boolean_expr)
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
				p.clause_boolean_expr(2)
			}

		}
		p.SetState(426)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())
	}

	return localctx
}

// ITag_cascade_exprContext is an interface to support dynamic dispatch.
type ITag_cascade_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTag_cascade_exprContext differentiates from other interfaces.
	IsTag_cascade_exprContext()
}

type Tag_cascade_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTag_cascade_exprContext() *Tag_cascade_exprContext {
	var p = new(Tag_cascade_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tag_cascade_expr
	return p
}

func (*Tag_cascade_exprContext) IsTag_cascade_exprContext() {}

func NewTag_cascade_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Tag_cascade_exprContext {
	var p = new(Tag_cascade_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tag_cascade_expr

	return p
}

func (s *Tag_cascade_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Tag_cascade_exprContext) Tag_equal_expr() ITag_equal_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_equal_exprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_equal_exprContext)
}

func (s *Tag_cascade_exprContext) Tag_boolean_expr() ITag_boolean_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_boolean_exprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_boolean_exprContext)
}

func (s *Tag_cascade_exprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *Tag_cascade_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Tag_cascade_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Tag_cascade_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTag_cascade_expr(s)
	}
}

func (s *Tag_cascade_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTag_cascade_expr(s)
	}
}

func (s *Tag_cascade_exprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTag_cascade_expr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Tag_cascade_expr() (localctx ITag_cascade_exprContext) {
	localctx = NewTag_cascade_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, SQLParserRULE_tag_cascade_expr)
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
			p.Tag_equal_expr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(428)
			p.tag_boolean_expr(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(429)
			p.Tag_equal_expr()
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
				p.tag_boolean_expr(0)
			}

		}

	}

	return localctx
}

// ITag_equal_exprContext is an interface to support dynamic dispatch.
type ITag_equal_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTag_equal_exprContext differentiates from other interfaces.
	IsTag_equal_exprContext()
}

type Tag_equal_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTag_equal_exprContext() *Tag_equal_exprContext {
	var p = new(Tag_equal_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tag_equal_expr
	return p
}

func (*Tag_equal_exprContext) IsTag_equal_exprContext() {}

func NewTag_equal_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Tag_equal_exprContext {
	var p = new(Tag_equal_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tag_equal_expr

	return p
}

func (s *Tag_equal_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Tag_equal_exprContext) T_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUE, 0)
}

func (s *Tag_equal_exprContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *Tag_equal_exprContext) Tag_value_pattern() ITag_value_patternContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_value_patternContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_value_patternContext)
}

func (s *Tag_equal_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Tag_equal_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Tag_equal_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTag_equal_expr(s)
	}
}

func (s *Tag_equal_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTag_equal_expr(s)
	}
}

func (s *Tag_equal_exprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTag_equal_expr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Tag_equal_expr() (localctx ITag_equal_exprContext) {
	localctx = NewTag_equal_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, SQLParserRULE_tag_equal_expr)

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
		p.Tag_value_pattern()
	}

	return localctx
}

// ITag_boolean_exprContext is an interface to support dynamic dispatch.
type ITag_boolean_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTag_boolean_exprContext differentiates from other interfaces.
	IsTag_boolean_exprContext()
}

type Tag_boolean_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTag_boolean_exprContext() *Tag_boolean_exprContext {
	var p = new(Tag_boolean_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tag_boolean_expr
	return p
}

func (*Tag_boolean_exprContext) IsTag_boolean_exprContext() {}

func NewTag_boolean_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Tag_boolean_exprContext {
	var p = new(Tag_boolean_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tag_boolean_expr

	return p
}

func (s *Tag_boolean_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Tag_boolean_exprContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *Tag_boolean_exprContext) AllTag_boolean_expr() []ITag_boolean_exprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITag_boolean_exprContext)(nil)).Elem())
	var tst = make([]ITag_boolean_exprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITag_boolean_exprContext)
		}
	}

	return tst
}

func (s *Tag_boolean_exprContext) Tag_boolean_expr(i int) ITag_boolean_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_boolean_exprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITag_boolean_exprContext)
}

func (s *Tag_boolean_exprContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *Tag_boolean_exprContext) Tag_key() ITag_keyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_keyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_keyContext)
}

func (s *Tag_boolean_exprContext) Tag_value() ITag_valueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_valueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_valueContext)
}

func (s *Tag_boolean_exprContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *Tag_boolean_exprContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *Tag_boolean_exprContext) T_REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REGEXP, 0)
}

func (s *Tag_boolean_exprContext) T_NOTEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL, 0)
}

func (s *Tag_boolean_exprContext) T_NOTEQUAL2() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL2, 0)
}

func (s *Tag_boolean_exprContext) Tag_value_list() ITag_value_listContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_value_listContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_value_listContext)
}

func (s *Tag_boolean_exprContext) T_IN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IN, 0)
}

func (s *Tag_boolean_exprContext) T_NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOT, 0)
}

func (s *Tag_boolean_exprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *Tag_boolean_exprContext) T_OR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OR, 0)
}

func (s *Tag_boolean_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Tag_boolean_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Tag_boolean_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTag_boolean_expr(s)
	}
}

func (s *Tag_boolean_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTag_boolean_expr(s)
	}
}

func (s *Tag_boolean_exprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTag_boolean_expr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Tag_boolean_expr() (localctx ITag_boolean_exprContext) {
	return p.tag_boolean_expr(0)
}

func (p *SQLParser) tag_boolean_expr(_p int) (localctx ITag_boolean_exprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewTag_boolean_exprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ITag_boolean_exprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 88
	p.EnterRecursionRule(localctx, 88, SQLParserRULE_tag_boolean_expr, _p)
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
			p.tag_boolean_expr(0)
		}
		{
			p.SetState(443)
			p.Match(SQLParserT_CLOSE_P)
		}

	case 2:
		{
			p.SetState(445)
			p.Tag_key()
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
			p.Tag_value()
		}

	case 3:
		{
			p.SetState(449)
			p.Tag_key()
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
			p.Tag_value_list()
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
			localctx = NewTag_boolean_exprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_tag_boolean_expr)
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
				p.tag_boolean_expr(2)
			}

		}
		p.SetState(468)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())
	}

	return localctx
}

// ITag_value_listContext is an interface to support dynamic dispatch.
type ITag_value_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTag_value_listContext differentiates from other interfaces.
	IsTag_value_listContext()
}

type Tag_value_listContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTag_value_listContext() *Tag_value_listContext {
	var p = new(Tag_value_listContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tag_value_list
	return p
}

func (*Tag_value_listContext) IsTag_value_listContext() {}

func NewTag_value_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Tag_value_listContext {
	var p = new(Tag_value_listContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tag_value_list

	return p
}

func (s *Tag_value_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Tag_value_listContext) AllTag_value() []ITag_valueContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITag_valueContext)(nil)).Elem())
	var tst = make([]ITag_valueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITag_valueContext)
		}
	}

	return tst
}

func (s *Tag_value_listContext) Tag_value(i int) ITag_valueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_valueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITag_valueContext)
}

func (s *Tag_value_listContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *Tag_value_listContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *Tag_value_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Tag_value_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Tag_value_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTag_value_list(s)
	}
}

func (s *Tag_value_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTag_value_list(s)
	}
}

func (s *Tag_value_listContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTag_value_list(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Tag_value_list() (localctx ITag_value_listContext) {
	localctx = NewTag_value_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, SQLParserRULE_tag_value_list)
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
		p.Tag_value()
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
			p.Tag_value()
		}

		p.SetState(476)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ITime_exprContext is an interface to support dynamic dispatch.
type ITime_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTime_exprContext differentiates from other interfaces.
	IsTime_exprContext()
}

type Time_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTime_exprContext() *Time_exprContext {
	var p = new(Time_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_time_expr
	return p
}

func (*Time_exprContext) IsTime_exprContext() {}

func NewTime_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Time_exprContext {
	var p = new(Time_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_time_expr

	return p
}

func (s *Time_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Time_exprContext) AllTime_boolean_expr() []ITime_boolean_exprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITime_boolean_exprContext)(nil)).Elem())
	var tst = make([]ITime_boolean_exprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITime_boolean_exprContext)
		}
	}

	return tst
}

func (s *Time_exprContext) Time_boolean_expr(i int) ITime_boolean_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITime_boolean_exprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITime_boolean_exprContext)
}

func (s *Time_exprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *Time_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Time_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Time_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTime_expr(s)
	}
}

func (s *Time_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTime_expr(s)
	}
}

func (s *Time_exprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTime_expr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Time_expr() (localctx ITime_exprContext) {
	localctx = NewTime_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, SQLParserRULE_time_expr)

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
		p.Time_boolean_expr()
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
			p.Time_boolean_expr()
		}

	}

	return localctx
}

// ITime_boolean_exprContext is an interface to support dynamic dispatch.
type ITime_boolean_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTime_boolean_exprContext differentiates from other interfaces.
	IsTime_boolean_exprContext()
}

type Time_boolean_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTime_boolean_exprContext() *Time_boolean_exprContext {
	var p = new(Time_boolean_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_time_boolean_expr
	return p
}

func (*Time_boolean_exprContext) IsTime_boolean_exprContext() {}

func NewTime_boolean_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Time_boolean_exprContext {
	var p = new(Time_boolean_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_time_boolean_expr

	return p
}

func (s *Time_boolean_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Time_boolean_exprContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *Time_boolean_exprContext) Bool_expr_binary_operator() IBool_expr_binary_operatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBool_expr_binary_operatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBool_expr_binary_operatorContext)
}

func (s *Time_boolean_exprContext) Now_expr() INow_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INow_exprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INow_exprContext)
}

func (s *Time_boolean_exprContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *Time_boolean_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Time_boolean_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Time_boolean_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTime_boolean_expr(s)
	}
}

func (s *Time_boolean_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTime_boolean_expr(s)
	}
}

func (s *Time_boolean_exprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTime_boolean_expr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Time_boolean_expr() (localctx ITime_boolean_exprContext) {
	localctx = NewTime_boolean_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, SQLParserRULE_time_boolean_expr)

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
		p.Bool_expr_binary_operator()
	}
	p.SetState(486)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_NOW:
		{
			p.SetState(484)
			p.Now_expr()
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

// INow_exprContext is an interface to support dynamic dispatch.
type INow_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNow_exprContext differentiates from other interfaces.
	IsNow_exprContext()
}

type Now_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNow_exprContext() *Now_exprContext {
	var p = new(Now_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_now_expr
	return p
}

func (*Now_exprContext) IsNow_exprContext() {}

func NewNow_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Now_exprContext {
	var p = new(Now_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_now_expr

	return p
}

func (s *Now_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Now_exprContext) Now_func() INow_funcContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INow_funcContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INow_funcContext)
}

func (s *Now_exprContext) Duration_lit() IDuration_litContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDuration_litContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDuration_litContext)
}

func (s *Now_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Now_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Now_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNow_expr(s)
	}
}

func (s *Now_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNow_expr(s)
	}
}

func (s *Now_exprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitNow_expr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Now_expr() (localctx INow_exprContext) {
	localctx = NewNow_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, SQLParserRULE_now_expr)

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
		p.Now_func()
	}
	p.SetState(490)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(489)
			p.Duration_lit()
		}

	}

	return localctx
}

// INow_funcContext is an interface to support dynamic dispatch.
type INow_funcContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNow_funcContext differentiates from other interfaces.
	IsNow_funcContext()
}

type Now_funcContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNow_funcContext() *Now_funcContext {
	var p = new(Now_funcContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_now_func
	return p
}

func (*Now_funcContext) IsNow_funcContext() {}

func NewNow_funcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Now_funcContext {
	var p = new(Now_funcContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_now_func

	return p
}

func (s *Now_funcContext) GetParser() antlr.Parser { return s.parser }

func (s *Now_funcContext) T_NOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOW, 0)
}

func (s *Now_funcContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *Now_funcContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *Now_funcContext) Expr_func_params() IExpr_func_paramsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpr_func_paramsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpr_func_paramsContext)
}

func (s *Now_funcContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Now_funcContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Now_funcContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNow_func(s)
	}
}

func (s *Now_funcContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNow_func(s)
	}
}

func (s *Now_funcContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitNow_func(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Now_func() (localctx INow_funcContext) {
	localctx = NewNow_funcContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, SQLParserRULE_now_func)
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
			p.Expr_func_params()
		}

	}
	{
		p.SetState(497)
		p.Match(SQLParserT_CLOSE_P)
	}

	return localctx
}

// IGroup_by_clauseContext is an interface to support dynamic dispatch.
type IGroup_by_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGroup_by_clauseContext differentiates from other interfaces.
	IsGroup_by_clauseContext()
}

type Group_by_clauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroup_by_clauseContext() *Group_by_clauseContext {
	var p = new(Group_by_clauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_group_by_clause
	return p
}

func (*Group_by_clauseContext) IsGroup_by_clauseContext() {}

func NewGroup_by_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Group_by_clauseContext {
	var p = new(Group_by_clauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_group_by_clause

	return p
}

func (s *Group_by_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *Group_by_clauseContext) T_GROUP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GROUP, 0)
}

func (s *Group_by_clauseContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *Group_by_clauseContext) Dimensions() IDimensionsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDimensionsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDimensionsContext)
}

func (s *Group_by_clauseContext) T_FILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FILL, 0)
}

func (s *Group_by_clauseContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *Group_by_clauseContext) Fill_option() IFill_optionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFill_optionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFill_optionContext)
}

func (s *Group_by_clauseContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *Group_by_clauseContext) Having_clause() IHaving_clauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHaving_clauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHaving_clauseContext)
}

func (s *Group_by_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Group_by_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Group_by_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterGroup_by_clause(s)
	}
}

func (s *Group_by_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitGroup_by_clause(s)
	}
}

func (s *Group_by_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitGroup_by_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Group_by_clause() (localctx IGroup_by_clauseContext) {
	localctx = NewGroup_by_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, SQLParserRULE_group_by_clause)
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
			p.Fill_option()
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
			p.Having_clause()
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

func (s *DimensionContext) Duration_lit() IDuration_litContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDuration_litContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDuration_litContext)
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
			p.Duration_lit()
		}
		{
			p.SetState(524)
			p.Match(SQLParserT_CLOSE_P)
		}

	}

	return localctx
}

// IFill_optionContext is an interface to support dynamic dispatch.
type IFill_optionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFill_optionContext differentiates from other interfaces.
	IsFill_optionContext()
}

type Fill_optionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFill_optionContext() *Fill_optionContext {
	var p = new(Fill_optionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_fill_option
	return p
}

func (*Fill_optionContext) IsFill_optionContext() {}

func NewFill_optionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Fill_optionContext {
	var p = new(Fill_optionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_fill_option

	return p
}

func (s *Fill_optionContext) GetParser() antlr.Parser { return s.parser }

func (s *Fill_optionContext) T_NULL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NULL, 0)
}

func (s *Fill_optionContext) T_PREVIOUS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PREVIOUS, 0)
}

func (s *Fill_optionContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *Fill_optionContext) L_DEC() antlr.TerminalNode {
	return s.GetToken(SQLParserL_DEC, 0)
}

func (s *Fill_optionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Fill_optionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Fill_optionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFill_option(s)
	}
}

func (s *Fill_optionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFill_option(s)
	}
}

func (s *Fill_optionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFill_option(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Fill_option() (localctx IFill_optionContext) {
	localctx = NewFill_optionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, SQLParserRULE_fill_option)
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

// IOrder_by_clauseContext is an interface to support dynamic dispatch.
type IOrder_by_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOrder_by_clauseContext differentiates from other interfaces.
	IsOrder_by_clauseContext()
}

type Order_by_clauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrder_by_clauseContext() *Order_by_clauseContext {
	var p = new(Order_by_clauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_order_by_clause
	return p
}

func (*Order_by_clauseContext) IsOrder_by_clauseContext() {}

func NewOrder_by_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Order_by_clauseContext {
	var p = new(Order_by_clauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_order_by_clause

	return p
}

func (s *Order_by_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *Order_by_clauseContext) T_ORDER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ORDER, 0)
}

func (s *Order_by_clauseContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *Order_by_clauseContext) Sort_fields() ISort_fieldsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISort_fieldsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISort_fieldsContext)
}

func (s *Order_by_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Order_by_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Order_by_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterOrder_by_clause(s)
	}
}

func (s *Order_by_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitOrder_by_clause(s)
	}
}

func (s *Order_by_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitOrder_by_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Order_by_clause() (localctx IOrder_by_clauseContext) {
	localctx = NewOrder_by_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, SQLParserRULE_order_by_clause)

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
		p.Sort_fields()
	}

	return localctx
}

// IInterval_by_clauseContext is an interface to support dynamic dispatch.
type IInterval_by_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInterval_by_clauseContext differentiates from other interfaces.
	IsInterval_by_clauseContext()
}

type Interval_by_clauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInterval_by_clauseContext() *Interval_by_clauseContext {
	var p = new(Interval_by_clauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_interval_by_clause
	return p
}

func (*Interval_by_clauseContext) IsInterval_by_clauseContext() {}

func NewInterval_by_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Interval_by_clauseContext {
	var p = new(Interval_by_clauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_interval_by_clause

	return p
}

func (s *Interval_by_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *Interval_by_clauseContext) T_INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL, 0)
}

func (s *Interval_by_clauseContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *Interval_by_clauseContext) Interval_name_val() IInterval_name_valContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInterval_name_valContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInterval_name_valContext)
}

func (s *Interval_by_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Interval_by_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Interval_by_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterInterval_by_clause(s)
	}
}

func (s *Interval_by_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitInterval_by_clause(s)
	}
}

func (s *Interval_by_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitInterval_by_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Interval_by_clause() (localctx IInterval_by_clauseContext) {
	localctx = NewInterval_by_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, SQLParserRULE_interval_by_clause)

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
		p.Interval_name_val()
	}

	return localctx
}

// ISort_fieldContext is an interface to support dynamic dispatch.
type ISort_fieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSort_fieldContext differentiates from other interfaces.
	IsSort_fieldContext()
}

type Sort_fieldContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySort_fieldContext() *Sort_fieldContext {
	var p = new(Sort_fieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_sort_field
	return p
}

func (*Sort_fieldContext) IsSort_fieldContext() {}

func NewSort_fieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Sort_fieldContext {
	var p = new(Sort_fieldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_sort_field

	return p
}

func (s *Sort_fieldContext) GetParser() antlr.Parser { return s.parser }

func (s *Sort_fieldContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *Sort_fieldContext) AllT_ASC() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_ASC)
}

func (s *Sort_fieldContext) T_ASC(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_ASC, i)
}

func (s *Sort_fieldContext) AllT_DESC() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_DESC)
}

func (s *Sort_fieldContext) T_DESC(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_DESC, i)
}

func (s *Sort_fieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Sort_fieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Sort_fieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSort_field(s)
	}
}

func (s *Sort_fieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSort_field(s)
	}
}

func (s *Sort_fieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitSort_field(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Sort_field() (localctx ISort_fieldContext) {
	localctx = NewSort_fieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, SQLParserRULE_sort_field)
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

// ISort_fieldsContext is an interface to support dynamic dispatch.
type ISort_fieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSort_fieldsContext differentiates from other interfaces.
	IsSort_fieldsContext()
}

type Sort_fieldsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySort_fieldsContext() *Sort_fieldsContext {
	var p = new(Sort_fieldsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_sort_fields
	return p
}

func (*Sort_fieldsContext) IsSort_fieldsContext() {}

func NewSort_fieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Sort_fieldsContext {
	var p = new(Sort_fieldsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_sort_fields

	return p
}

func (s *Sort_fieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *Sort_fieldsContext) AllSort_field() []ISort_fieldContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISort_fieldContext)(nil)).Elem())
	var tst = make([]ISort_fieldContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISort_fieldContext)
		}
	}

	return tst
}

func (s *Sort_fieldsContext) Sort_field(i int) ISort_fieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISort_fieldContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISort_fieldContext)
}

func (s *Sort_fieldsContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *Sort_fieldsContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *Sort_fieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Sort_fieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Sort_fieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSort_fields(s)
	}
}

func (s *Sort_fieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSort_fields(s)
	}
}

func (s *Sort_fieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitSort_fields(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Sort_fields() (localctx ISort_fieldsContext) {
	localctx = NewSort_fieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, SQLParserRULE_sort_fields)
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
		p.Sort_field()
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
			p.Sort_field()
		}

		p.SetState(552)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IHaving_clauseContext is an interface to support dynamic dispatch.
type IHaving_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHaving_clauseContext differentiates from other interfaces.
	IsHaving_clauseContext()
}

type Having_clauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHaving_clauseContext() *Having_clauseContext {
	var p = new(Having_clauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_having_clause
	return p
}

func (*Having_clauseContext) IsHaving_clauseContext() {}

func NewHaving_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Having_clauseContext {
	var p = new(Having_clauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_having_clause

	return p
}

func (s *Having_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *Having_clauseContext) T_HAVING() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HAVING, 0)
}

func (s *Having_clauseContext) Bool_expr() IBool_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBool_exprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBool_exprContext)
}

func (s *Having_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Having_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Having_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterHaving_clause(s)
	}
}

func (s *Having_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitHaving_clause(s)
	}
}

func (s *Having_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitHaving_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Having_clause() (localctx IHaving_clauseContext) {
	localctx = NewHaving_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, SQLParserRULE_having_clause)

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
		p.bool_expr(0)
	}

	return localctx
}

// IBool_exprContext is an interface to support dynamic dispatch.
type IBool_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBool_exprContext differentiates from other interfaces.
	IsBool_exprContext()
}

type Bool_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBool_exprContext() *Bool_exprContext {
	var p = new(Bool_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_bool_expr
	return p
}

func (*Bool_exprContext) IsBool_exprContext() {}

func NewBool_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Bool_exprContext {
	var p = new(Bool_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_bool_expr

	return p
}

func (s *Bool_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Bool_exprContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *Bool_exprContext) AllBool_expr() []IBool_exprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IBool_exprContext)(nil)).Elem())
	var tst = make([]IBool_exprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IBool_exprContext)
		}
	}

	return tst
}

func (s *Bool_exprContext) Bool_expr(i int) IBool_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBool_exprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IBool_exprContext)
}

func (s *Bool_exprContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *Bool_exprContext) Bool_expr_atom() IBool_expr_atomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBool_expr_atomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBool_expr_atomContext)
}

func (s *Bool_exprContext) Bool_expr_logical_op() IBool_expr_logical_opContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBool_expr_logical_opContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBool_expr_logical_opContext)
}

func (s *Bool_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Bool_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Bool_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBool_expr(s)
	}
}

func (s *Bool_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBool_expr(s)
	}
}

func (s *Bool_exprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBool_expr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Bool_expr() (localctx IBool_exprContext) {
	return p.bool_expr(0)
}

func (p *SQLParser) bool_expr(_p int) (localctx IBool_exprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewBool_exprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IBool_exprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 118
	p.EnterRecursionRule(localctx, 118, SQLParserRULE_bool_expr, _p)

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
			p.bool_expr(0)
		}
		{
			p.SetState(559)
			p.Match(SQLParserT_CLOSE_P)
		}

	case 2:
		{
			p.SetState(561)
			p.Bool_expr_atom()
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
			localctx = NewBool_exprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_bool_expr)
			p.SetState(564)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(565)
				p.Bool_expr_logical_op()
			}
			{
				p.SetState(566)
				p.bool_expr(3)
			}

		}
		p.SetState(572)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext())
	}

	return localctx
}

// IBool_expr_logical_opContext is an interface to support dynamic dispatch.
type IBool_expr_logical_opContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBool_expr_logical_opContext differentiates from other interfaces.
	IsBool_expr_logical_opContext()
}

type Bool_expr_logical_opContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBool_expr_logical_opContext() *Bool_expr_logical_opContext {
	var p = new(Bool_expr_logical_opContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_bool_expr_logical_op
	return p
}

func (*Bool_expr_logical_opContext) IsBool_expr_logical_opContext() {}

func NewBool_expr_logical_opContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Bool_expr_logical_opContext {
	var p = new(Bool_expr_logical_opContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_bool_expr_logical_op

	return p
}

func (s *Bool_expr_logical_opContext) GetParser() antlr.Parser { return s.parser }

func (s *Bool_expr_logical_opContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *Bool_expr_logical_opContext) T_OR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OR, 0)
}

func (s *Bool_expr_logical_opContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Bool_expr_logical_opContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Bool_expr_logical_opContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBool_expr_logical_op(s)
	}
}

func (s *Bool_expr_logical_opContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBool_expr_logical_op(s)
	}
}

func (s *Bool_expr_logical_opContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBool_expr_logical_op(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Bool_expr_logical_op() (localctx IBool_expr_logical_opContext) {
	localctx = NewBool_expr_logical_opContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 120, SQLParserRULE_bool_expr_logical_op)
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

// IBool_expr_atomContext is an interface to support dynamic dispatch.
type IBool_expr_atomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBool_expr_atomContext differentiates from other interfaces.
	IsBool_expr_atomContext()
}

type Bool_expr_atomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBool_expr_atomContext() *Bool_expr_atomContext {
	var p = new(Bool_expr_atomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_bool_expr_atom
	return p
}

func (*Bool_expr_atomContext) IsBool_expr_atomContext() {}

func NewBool_expr_atomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Bool_expr_atomContext {
	var p = new(Bool_expr_atomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_bool_expr_atom

	return p
}

func (s *Bool_expr_atomContext) GetParser() antlr.Parser { return s.parser }

func (s *Bool_expr_atomContext) Bool_expr_binary() IBool_expr_binaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBool_expr_binaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBool_expr_binaryContext)
}

func (s *Bool_expr_atomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Bool_expr_atomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Bool_expr_atomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBool_expr_atom(s)
	}
}

func (s *Bool_expr_atomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBool_expr_atom(s)
	}
}

func (s *Bool_expr_atomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBool_expr_atom(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Bool_expr_atom() (localctx IBool_expr_atomContext) {
	localctx = NewBool_expr_atomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, SQLParserRULE_bool_expr_atom)

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
		p.Bool_expr_binary()
	}

	return localctx
}

// IBool_expr_binaryContext is an interface to support dynamic dispatch.
type IBool_expr_binaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBool_expr_binaryContext differentiates from other interfaces.
	IsBool_expr_binaryContext()
}

type Bool_expr_binaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBool_expr_binaryContext() *Bool_expr_binaryContext {
	var p = new(Bool_expr_binaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_bool_expr_binary
	return p
}

func (*Bool_expr_binaryContext) IsBool_expr_binaryContext() {}

func NewBool_expr_binaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Bool_expr_binaryContext {
	var p = new(Bool_expr_binaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_bool_expr_binary

	return p
}

func (s *Bool_expr_binaryContext) GetParser() antlr.Parser { return s.parser }

func (s *Bool_expr_binaryContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *Bool_expr_binaryContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *Bool_expr_binaryContext) Bool_expr_binary_operator() IBool_expr_binary_operatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBool_expr_binary_operatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBool_expr_binary_operatorContext)
}

func (s *Bool_expr_binaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Bool_expr_binaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Bool_expr_binaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBool_expr_binary(s)
	}
}

func (s *Bool_expr_binaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBool_expr_binary(s)
	}
}

func (s *Bool_expr_binaryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBool_expr_binary(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Bool_expr_binary() (localctx IBool_expr_binaryContext) {
	localctx = NewBool_expr_binaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, SQLParserRULE_bool_expr_binary)

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
		p.Bool_expr_binary_operator()
	}
	{
		p.SetState(579)
		p.expr(0)
	}

	return localctx
}

// IBool_expr_binary_operatorContext is an interface to support dynamic dispatch.
type IBool_expr_binary_operatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBool_expr_binary_operatorContext differentiates from other interfaces.
	IsBool_expr_binary_operatorContext()
}

type Bool_expr_binary_operatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBool_expr_binary_operatorContext() *Bool_expr_binary_operatorContext {
	var p = new(Bool_expr_binary_operatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_bool_expr_binary_operator
	return p
}

func (*Bool_expr_binary_operatorContext) IsBool_expr_binary_operatorContext() {}

func NewBool_expr_binary_operatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Bool_expr_binary_operatorContext {
	var p = new(Bool_expr_binary_operatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_bool_expr_binary_operator

	return p
}

func (s *Bool_expr_binary_operatorContext) GetParser() antlr.Parser { return s.parser }

func (s *Bool_expr_binary_operatorContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *Bool_expr_binary_operatorContext) T_NOTEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL, 0)
}

func (s *Bool_expr_binary_operatorContext) T_NOTEQUAL2() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL2, 0)
}

func (s *Bool_expr_binary_operatorContext) T_LESS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LESS, 0)
}

func (s *Bool_expr_binary_operatorContext) T_LESSEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LESSEQUAL, 0)
}

func (s *Bool_expr_binary_operatorContext) T_GREATER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GREATER, 0)
}

func (s *Bool_expr_binary_operatorContext) T_GREATEREQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GREATEREQUAL, 0)
}

func (s *Bool_expr_binary_operatorContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *Bool_expr_binary_operatorContext) T_REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REGEXP, 0)
}

func (s *Bool_expr_binary_operatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Bool_expr_binary_operatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Bool_expr_binary_operatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBool_expr_binary_operator(s)
	}
}

func (s *Bool_expr_binary_operatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBool_expr_binary_operator(s)
	}
}

func (s *Bool_expr_binary_operatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBool_expr_binary_operator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Bool_expr_binary_operator() (localctx IBool_expr_binary_operatorContext) {
	localctx = NewBool_expr_binary_operatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, SQLParserRULE_bool_expr_binary_operator)
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

func (s *ExprContext) Expr_func() IExpr_funcContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpr_funcContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpr_funcContext)
}

func (s *ExprContext) Expr_atom() IExpr_atomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpr_atomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpr_atomContext)
}

func (s *ExprContext) Duration_lit() IDuration_litContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDuration_litContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDuration_litContext)
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
			p.Expr_func()
		}

	case 3:
		{
			p.SetState(597)
			p.Expr_atom()
		}

	case 4:
		{
			p.SetState(598)
			p.Duration_lit()
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

// IDuration_litContext is an interface to support dynamic dispatch.
type IDuration_litContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDuration_litContext differentiates from other interfaces.
	IsDuration_litContext()
}

type Duration_litContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDuration_litContext() *Duration_litContext {
	var p = new(Duration_litContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_duration_lit
	return p
}

func (*Duration_litContext) IsDuration_litContext() {}

func NewDuration_litContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Duration_litContext {
	var p = new(Duration_litContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_duration_lit

	return p
}

func (s *Duration_litContext) GetParser() antlr.Parser { return s.parser }

func (s *Duration_litContext) Int_number() IInt_numberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInt_numberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInt_numberContext)
}

func (s *Duration_litContext) Interval_item() IInterval_itemContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInterval_itemContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInterval_itemContext)
}

func (s *Duration_litContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Duration_litContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Duration_litContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDuration_lit(s)
	}
}

func (s *Duration_litContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDuration_lit(s)
	}
}

func (s *Duration_litContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDuration_lit(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Duration_lit() (localctx IDuration_litContext) {
	localctx = NewDuration_litContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 130, SQLParserRULE_duration_lit)

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
		p.Int_number()
	}
	{
		p.SetState(619)
		p.Interval_item()
	}

	return localctx
}

// IInterval_itemContext is an interface to support dynamic dispatch.
type IInterval_itemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInterval_itemContext differentiates from other interfaces.
	IsInterval_itemContext()
}

type Interval_itemContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInterval_itemContext() *Interval_itemContext {
	var p = new(Interval_itemContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_interval_item
	return p
}

func (*Interval_itemContext) IsInterval_itemContext() {}

func NewInterval_itemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Interval_itemContext {
	var p = new(Interval_itemContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_interval_item

	return p
}

func (s *Interval_itemContext) GetParser() antlr.Parser { return s.parser }

func (s *Interval_itemContext) T_SECOND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SECOND, 0)
}

func (s *Interval_itemContext) T_MINUTE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MINUTE, 0)
}

func (s *Interval_itemContext) T_HOUR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HOUR, 0)
}

func (s *Interval_itemContext) T_DAY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DAY, 0)
}

func (s *Interval_itemContext) T_WEEK() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WEEK, 0)
}

func (s *Interval_itemContext) T_MONTH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MONTH, 0)
}

func (s *Interval_itemContext) T_YEAR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_YEAR, 0)
}

func (s *Interval_itemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Interval_itemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Interval_itemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterInterval_item(s)
	}
}

func (s *Interval_itemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitInterval_item(s)
	}
}

func (s *Interval_itemContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitInterval_item(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Interval_item() (localctx IInterval_itemContext) {
	localctx = NewInterval_itemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 132, SQLParserRULE_interval_item)
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

// IExpr_funcContext is an interface to support dynamic dispatch.
type IExpr_funcContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpr_funcContext differentiates from other interfaces.
	IsExpr_funcContext()
}

type Expr_funcContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpr_funcContext() *Expr_funcContext {
	var p = new(Expr_funcContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_expr_func
	return p
}

func (*Expr_funcContext) IsExpr_funcContext() {}

func NewExpr_funcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Expr_funcContext {
	var p = new(Expr_funcContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_expr_func

	return p
}

func (s *Expr_funcContext) GetParser() antlr.Parser { return s.parser }

func (s *Expr_funcContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *Expr_funcContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *Expr_funcContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *Expr_funcContext) Expr_func_params() IExpr_func_paramsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpr_func_paramsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpr_func_paramsContext)
}

func (s *Expr_funcContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Expr_funcContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Expr_funcContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterExpr_func(s)
	}
}

func (s *Expr_funcContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitExpr_func(s)
	}
}

func (s *Expr_funcContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitExpr_func(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Expr_func() (localctx IExpr_funcContext) {
	localctx = NewExpr_funcContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 134, SQLParserRULE_expr_func)
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
			p.Expr_func_params()
		}

	}
	{
		p.SetState(628)
		p.Match(SQLParserT_CLOSE_P)
	}

	return localctx
}

// IExpr_func_paramsContext is an interface to support dynamic dispatch.
type IExpr_func_paramsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpr_func_paramsContext differentiates from other interfaces.
	IsExpr_func_paramsContext()
}

type Expr_func_paramsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpr_func_paramsContext() *Expr_func_paramsContext {
	var p = new(Expr_func_paramsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_expr_func_params
	return p
}

func (*Expr_func_paramsContext) IsExpr_func_paramsContext() {}

func NewExpr_func_paramsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Expr_func_paramsContext {
	var p = new(Expr_func_paramsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_expr_func_params

	return p
}

func (s *Expr_func_paramsContext) GetParser() antlr.Parser { return s.parser }

func (s *Expr_func_paramsContext) AllFunc_param() []IFunc_paramContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFunc_paramContext)(nil)).Elem())
	var tst = make([]IFunc_paramContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFunc_paramContext)
		}
	}

	return tst
}

func (s *Expr_func_paramsContext) Func_param(i int) IFunc_paramContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFunc_paramContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFunc_paramContext)
}

func (s *Expr_func_paramsContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *Expr_func_paramsContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *Expr_func_paramsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Expr_func_paramsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Expr_func_paramsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterExpr_func_params(s)
	}
}

func (s *Expr_func_paramsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitExpr_func_params(s)
	}
}

func (s *Expr_func_paramsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitExpr_func_params(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Expr_func_params() (localctx IExpr_func_paramsContext) {
	localctx = NewExpr_func_paramsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 136, SQLParserRULE_expr_func_params)
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
		p.Func_param()
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
			p.Func_param()
		}

		p.SetState(637)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IFunc_paramContext is an interface to support dynamic dispatch.
type IFunc_paramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunc_paramContext differentiates from other interfaces.
	IsFunc_paramContext()
}

type Func_paramContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunc_paramContext() *Func_paramContext {
	var p = new(Func_paramContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_func_param
	return p
}

func (*Func_paramContext) IsFunc_paramContext() {}

func NewFunc_paramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Func_paramContext {
	var p = new(Func_paramContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_func_param

	return p
}

func (s *Func_paramContext) GetParser() antlr.Parser { return s.parser }

func (s *Func_paramContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *Func_paramContext) Tag_boolean_expr() ITag_boolean_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_boolean_exprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_boolean_exprContext)
}

func (s *Func_paramContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Func_paramContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Func_paramContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFunc_param(s)
	}
}

func (s *Func_paramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFunc_param(s)
	}
}

func (s *Func_paramContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFunc_param(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Func_param() (localctx IFunc_paramContext) {
	localctx = NewFunc_paramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 138, SQLParserRULE_func_param)

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
			p.tag_boolean_expr(0)
		}

	}

	return localctx
}

// IExpr_atomContext is an interface to support dynamic dispatch.
type IExpr_atomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpr_atomContext differentiates from other interfaces.
	IsExpr_atomContext()
}

type Expr_atomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpr_atomContext() *Expr_atomContext {
	var p = new(Expr_atomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_expr_atom
	return p
}

func (*Expr_atomContext) IsExpr_atomContext() {}

func NewExpr_atomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Expr_atomContext {
	var p = new(Expr_atomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_expr_atom

	return p
}

func (s *Expr_atomContext) GetParser() antlr.Parser { return s.parser }

func (s *Expr_atomContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *Expr_atomContext) Ident_filter() IIdent_filterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdent_filterContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdent_filterContext)
}

func (s *Expr_atomContext) Dec_number() IDec_numberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDec_numberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDec_numberContext)
}

func (s *Expr_atomContext) Int_number() IInt_numberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInt_numberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInt_numberContext)
}

func (s *Expr_atomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Expr_atomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Expr_atomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterExpr_atom(s)
	}
}

func (s *Expr_atomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitExpr_atom(s)
	}
}

func (s *Expr_atomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitExpr_atom(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Expr_atom() (localctx IExpr_atomContext) {
	localctx = NewExpr_atomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 140, SQLParserRULE_expr_atom)

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
				p.Ident_filter()
			}

		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(646)
			p.Dec_number()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(647)
			p.Int_number()
		}

	}

	return localctx
}

// IIdent_filterContext is an interface to support dynamic dispatch.
type IIdent_filterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIdent_filterContext differentiates from other interfaces.
	IsIdent_filterContext()
}

type Ident_filterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdent_filterContext() *Ident_filterContext {
	var p = new(Ident_filterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_ident_filter
	return p
}

func (*Ident_filterContext) IsIdent_filterContext() {}

func NewIdent_filterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Ident_filterContext {
	var p = new(Ident_filterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_ident_filter

	return p
}

func (s *Ident_filterContext) GetParser() antlr.Parser { return s.parser }

func (s *Ident_filterContext) T_OPEN_SB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_SB, 0)
}

func (s *Ident_filterContext) Tag_boolean_expr() ITag_boolean_exprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITag_boolean_exprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITag_boolean_exprContext)
}

func (s *Ident_filterContext) T_CLOSE_SB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_SB, 0)
}

func (s *Ident_filterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Ident_filterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Ident_filterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIdent_filter(s)
	}
}

func (s *Ident_filterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIdent_filter(s)
	}
}

func (s *Ident_filterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIdent_filter(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Ident_filter() (localctx IIdent_filterContext) {
	localctx = NewIdent_filterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 142, SQLParserRULE_ident_filter)

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
		p.tag_boolean_expr(0)
	}
	{
		p.SetState(652)
		p.Match(SQLParserT_CLOSE_SB)
	}

	return localctx
}

// IInt_numberContext is an interface to support dynamic dispatch.
type IInt_numberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInt_numberContext differentiates from other interfaces.
	IsInt_numberContext()
}

type Int_numberContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInt_numberContext() *Int_numberContext {
	var p = new(Int_numberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_int_number
	return p
}

func (*Int_numberContext) IsInt_numberContext() {}

func NewInt_numberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Int_numberContext {
	var p = new(Int_numberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_int_number

	return p
}

func (s *Int_numberContext) GetParser() antlr.Parser { return s.parser }

func (s *Int_numberContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *Int_numberContext) T_SUB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUB, 0)
}

func (s *Int_numberContext) T_ADD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ADD, 0)
}

func (s *Int_numberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Int_numberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Int_numberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterInt_number(s)
	}
}

func (s *Int_numberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitInt_number(s)
	}
}

func (s *Int_numberContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitInt_number(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Int_number() (localctx IInt_numberContext) {
	localctx = NewInt_numberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 144, SQLParserRULE_int_number)
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

// IDec_numberContext is an interface to support dynamic dispatch.
type IDec_numberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDec_numberContext differentiates from other interfaces.
	IsDec_numberContext()
}

type Dec_numberContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDec_numberContext() *Dec_numberContext {
	var p = new(Dec_numberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_dec_number
	return p
}

func (*Dec_numberContext) IsDec_numberContext() {}

func NewDec_numberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Dec_numberContext {
	var p = new(Dec_numberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_dec_number

	return p
}

func (s *Dec_numberContext) GetParser() antlr.Parser { return s.parser }

func (s *Dec_numberContext) L_DEC() antlr.TerminalNode {
	return s.GetToken(SQLParserL_DEC, 0)
}

func (s *Dec_numberContext) T_SUB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUB, 0)
}

func (s *Dec_numberContext) T_ADD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ADD, 0)
}

func (s *Dec_numberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Dec_numberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Dec_numberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDec_number(s)
	}
}

func (s *Dec_numberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDec_number(s)
	}
}

func (s *Dec_numberContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDec_number(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Dec_number() (localctx IDec_numberContext) {
	localctx = NewDec_numberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 146, SQLParserRULE_dec_number)
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

// ILimit_clauseContext is an interface to support dynamic dispatch.
type ILimit_clauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLimit_clauseContext differentiates from other interfaces.
	IsLimit_clauseContext()
}

type Limit_clauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLimit_clauseContext() *Limit_clauseContext {
	var p = new(Limit_clauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_limit_clause
	return p
}

func (*Limit_clauseContext) IsLimit_clauseContext() {}

func NewLimit_clauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Limit_clauseContext {
	var p = new(Limit_clauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_limit_clause

	return p
}

func (s *Limit_clauseContext) GetParser() antlr.Parser { return s.parser }

func (s *Limit_clauseContext) T_LIMIT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIMIT, 0)
}

func (s *Limit_clauseContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *Limit_clauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Limit_clauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Limit_clauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterLimit_clause(s)
	}
}

func (s *Limit_clauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitLimit_clause(s)
	}
}

func (s *Limit_clauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitLimit_clause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Limit_clause() (localctx ILimit_clauseContext) {
	localctx = NewLimit_clauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 148, SQLParserRULE_limit_clause)

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

// IMetric_nameContext is an interface to support dynamic dispatch.
type IMetric_nameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMetric_nameContext differentiates from other interfaces.
	IsMetric_nameContext()
}

type Metric_nameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetric_nameContext() *Metric_nameContext {
	var p = new(Metric_nameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_metric_name
	return p
}

func (*Metric_nameContext) IsMetric_nameContext() {}

func NewMetric_nameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Metric_nameContext {
	var p = new(Metric_nameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_metric_name

	return p
}

func (s *Metric_nameContext) GetParser() antlr.Parser { return s.parser }

func (s *Metric_nameContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *Metric_nameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Metric_nameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Metric_nameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterMetric_name(s)
	}
}

func (s *Metric_nameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitMetric_name(s)
	}
}

func (s *Metric_nameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitMetric_name(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Metric_name() (localctx IMetric_nameContext) {
	localctx = NewMetric_nameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 150, SQLParserRULE_metric_name)

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

// ITag_keyContext is an interface to support dynamic dispatch.
type ITag_keyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTag_keyContext differentiates from other interfaces.
	IsTag_keyContext()
}

type Tag_keyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTag_keyContext() *Tag_keyContext {
	var p = new(Tag_keyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tag_key
	return p
}

func (*Tag_keyContext) IsTag_keyContext() {}

func NewTag_keyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Tag_keyContext {
	var p = new(Tag_keyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tag_key

	return p
}

func (s *Tag_keyContext) GetParser() antlr.Parser { return s.parser }

func (s *Tag_keyContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *Tag_keyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Tag_keyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Tag_keyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTag_key(s)
	}
}

func (s *Tag_keyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTag_key(s)
	}
}

func (s *Tag_keyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTag_key(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Tag_key() (localctx ITag_keyContext) {
	localctx = NewTag_keyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 152, SQLParserRULE_tag_key)

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

// ITag_valueContext is an interface to support dynamic dispatch.
type ITag_valueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTag_valueContext differentiates from other interfaces.
	IsTag_valueContext()
}

type Tag_valueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTag_valueContext() *Tag_valueContext {
	var p = new(Tag_valueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tag_value
	return p
}

func (*Tag_valueContext) IsTag_valueContext() {}

func NewTag_valueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Tag_valueContext {
	var p = new(Tag_valueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tag_value

	return p
}

func (s *Tag_valueContext) GetParser() antlr.Parser { return s.parser }

func (s *Tag_valueContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *Tag_valueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Tag_valueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Tag_valueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTag_value(s)
	}
}

func (s *Tag_valueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTag_value(s)
	}
}

func (s *Tag_valueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTag_value(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Tag_value() (localctx ITag_valueContext) {
	localctx = NewTag_valueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 154, SQLParserRULE_tag_value)

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

// ITag_value_patternContext is an interface to support dynamic dispatch.
type ITag_value_patternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTag_value_patternContext differentiates from other interfaces.
	IsTag_value_patternContext()
}

type Tag_value_patternContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTag_value_patternContext() *Tag_value_patternContext {
	var p = new(Tag_value_patternContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tag_value_pattern
	return p
}

func (*Tag_value_patternContext) IsTag_value_patternContext() {}

func NewTag_value_patternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Tag_value_patternContext {
	var p = new(Tag_value_patternContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tag_value_pattern

	return p
}

func (s *Tag_value_patternContext) GetParser() antlr.Parser { return s.parser }

func (s *Tag_value_patternContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *Tag_value_patternContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Tag_value_patternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Tag_value_patternContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTag_value_pattern(s)
	}
}

func (s *Tag_value_patternContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTag_value_pattern(s)
	}
}

func (s *Tag_value_patternContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTag_value_pattern(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Tag_value_pattern() (localctx ITag_value_patternContext) {
	localctx = NewTag_value_patternContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 156, SQLParserRULE_tag_value_pattern)

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

func (s *IdentContext) AllNon_reserved_words() []INon_reserved_wordsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*INon_reserved_wordsContext)(nil)).Elem())
	var tst = make([]INon_reserved_wordsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(INon_reserved_wordsContext)
		}
	}

	return tst
}

func (s *IdentContext) Non_reserved_words(i int) INon_reserved_wordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INon_reserved_wordsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(INon_reserved_wordsContext)
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
			p.Non_reserved_words()
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
					p.Non_reserved_words()
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

// INon_reserved_wordsContext is an interface to support dynamic dispatch.
type INon_reserved_wordsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNon_reserved_wordsContext differentiates from other interfaces.
	IsNon_reserved_wordsContext()
}

type Non_reserved_wordsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNon_reserved_wordsContext() *Non_reserved_wordsContext {
	var p = new(Non_reserved_wordsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_non_reserved_words
	return p
}

func (*Non_reserved_wordsContext) IsNon_reserved_wordsContext() {}

func NewNon_reserved_wordsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Non_reserved_wordsContext {
	var p = new(Non_reserved_wordsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_non_reserved_words

	return p
}

func (s *Non_reserved_wordsContext) GetParser() antlr.Parser { return s.parser }

func (s *Non_reserved_wordsContext) T_CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CREATE, 0)
}

func (s *Non_reserved_wordsContext) T_INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL, 0)
}

func (s *Non_reserved_wordsContext) T_SHARD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHARD, 0)
}

func (s *Non_reserved_wordsContext) T_REPLICATION() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REPLICATION, 0)
}

func (s *Non_reserved_wordsContext) T_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TTL, 0)
}

func (s *Non_reserved_wordsContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *Non_reserved_wordsContext) T_KILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KILL, 0)
}

func (s *Non_reserved_wordsContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *Non_reserved_wordsContext) T_DATASBAES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAES, 0)
}

func (s *Non_reserved_wordsContext) T_NODE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NODE, 0)
}

func (s *Non_reserved_wordsContext) T_MEASUREMENTS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MEASUREMENTS, 0)
}

func (s *Non_reserved_wordsContext) T_MEASUREMENT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MEASUREMENT, 0)
}

func (s *Non_reserved_wordsContext) T_FIELD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIELD, 0)
}

func (s *Non_reserved_wordsContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *Non_reserved_wordsContext) T_KEYS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEYS, 0)
}

func (s *Non_reserved_wordsContext) T_KEY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEY, 0)
}

func (s *Non_reserved_wordsContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *Non_reserved_wordsContext) T_VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUES, 0)
}

func (s *Non_reserved_wordsContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *Non_reserved_wordsContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *Non_reserved_wordsContext) T_LIMIT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIMIT, 0)
}

func (s *Non_reserved_wordsContext) T_QUERIES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUERIES, 0)
}

func (s *Non_reserved_wordsContext) T_QUERY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUERY, 0)
}

func (s *Non_reserved_wordsContext) T_SELECT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SELECT, 0)
}

func (s *Non_reserved_wordsContext) T_AS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AS, 0)
}

func (s *Non_reserved_wordsContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *Non_reserved_wordsContext) T_OR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OR, 0)
}

func (s *Non_reserved_wordsContext) T_NULL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NULL, 0)
}

func (s *Non_reserved_wordsContext) T_PREVIOUS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PREVIOUS, 0)
}

func (s *Non_reserved_wordsContext) T_FILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FILL, 0)
}

func (s *Non_reserved_wordsContext) T_ORDER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ORDER, 0)
}

func (s *Non_reserved_wordsContext) T_ASC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ASC, 0)
}

func (s *Non_reserved_wordsContext) T_DESC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DESC, 0)
}

func (s *Non_reserved_wordsContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *Non_reserved_wordsContext) T_NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOT, 0)
}

func (s *Non_reserved_wordsContext) T_BETWEEN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BETWEEN, 0)
}

func (s *Non_reserved_wordsContext) T_IS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IS, 0)
}

func (s *Non_reserved_wordsContext) T_PROFILE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PROFILE, 0)
}

func (s *Non_reserved_wordsContext) T_GROUP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GROUP, 0)
}

func (s *Non_reserved_wordsContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *Non_reserved_wordsContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *Non_reserved_wordsContext) T_STATS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATS, 0)
}

func (s *Non_reserved_wordsContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *Non_reserved_wordsContext) T_FOR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FOR, 0)
}

func (s *Non_reserved_wordsContext) T_SECOND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SECOND, 0)
}

func (s *Non_reserved_wordsContext) T_MINUTE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MINUTE, 0)
}

func (s *Non_reserved_wordsContext) T_HOUR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HOUR, 0)
}

func (s *Non_reserved_wordsContext) T_DAY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DAY, 0)
}

func (s *Non_reserved_wordsContext) T_WEEK() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WEEK, 0)
}

func (s *Non_reserved_wordsContext) T_MONTH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MONTH, 0)
}

func (s *Non_reserved_wordsContext) T_YEAR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_YEAR, 0)
}

func (s *Non_reserved_wordsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Non_reserved_wordsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Non_reserved_wordsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNon_reserved_words(s)
	}
}

func (s *Non_reserved_wordsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNon_reserved_words(s)
	}
}

func (s *Non_reserved_wordsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitNon_reserved_words(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Non_reserved_words() (localctx INon_reserved_wordsContext) {
	localctx = NewNon_reserved_wordsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 160, SQLParserRULE_non_reserved_words)
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
		var t *Clause_boolean_exprContext = nil
		if localctx != nil {
			t = localctx.(*Clause_boolean_exprContext)
		}
		return p.Clause_boolean_expr_Sempred(t, predIndex)

	case 44:
		var t *Tag_boolean_exprContext = nil
		if localctx != nil {
			t = localctx.(*Tag_boolean_exprContext)
		}
		return p.Tag_boolean_expr_Sempred(t, predIndex)

	case 59:
		var t *Bool_exprContext = nil
		if localctx != nil {
			t = localctx.(*Bool_exprContext)
		}
		return p.Bool_expr_Sempred(t, predIndex)

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

func (p *SQLParser) Clause_boolean_expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) Tag_boolean_expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 1:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) Bool_expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
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
