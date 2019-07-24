// Code generated from /Users/huangjie/Documents/github/lindb/sql/grammar/SQL.g4 by ANTLR 4.7.2. DO NOT EDIT.

package grammar // SQL
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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 95, 400, 
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7, 
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13, 
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9, 
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23, 
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4, 
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34, 
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9, 
	39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44, 9, 44, 
	4, 45, 9, 45, 4, 46, 9, 46, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 4, 5, 4, 99, 
	10, 4, 3, 4, 3, 4, 3, 4, 5, 4, 104, 10, 4, 3, 4, 5, 4, 107, 10, 4, 3, 4, 
	5, 4, 110, 10, 4, 3, 4, 5, 4, 113, 10, 4, 3, 4, 5, 4, 116, 10, 4, 3, 5, 
	3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 7, 6, 124, 10, 6, 12, 6, 14, 6, 127, 11, 
	6, 3, 7, 3, 7, 5, 7, 131, 10, 7, 3, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 
	10, 3, 10, 3, 10, 3, 11, 3, 11, 5, 11, 144, 10, 11, 3, 12, 3, 12, 3, 12, 
	3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 
	12, 5, 12, 160, 10, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 5, 12, 
	168, 10, 12, 3, 12, 3, 12, 3, 12, 3, 12, 5, 12, 174, 10, 12, 3, 12, 3, 
	12, 3, 12, 7, 12, 179, 10, 12, 12, 12, 14, 12, 182, 11, 12, 3, 13, 3, 13, 
	3, 13, 7, 13, 187, 10, 13, 12, 13, 14, 13, 190, 11, 13, 3, 14, 3, 14, 3, 
	14, 5, 14, 195, 10, 14, 3, 15, 3, 15, 3, 15, 3, 15, 5, 15, 201, 10, 15, 
	3, 16, 3, 16, 5, 16, 205, 10, 16, 3, 17, 3, 17, 3, 17, 5, 17, 210, 10, 
	17, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 3, 18, 3, 18, 3, 18, 3, 18, 3, 18, 
	5, 18, 222, 10, 18, 3, 18, 5, 18, 225, 10, 18, 3, 19, 3, 19, 3, 19, 7, 
	19, 230, 10, 19, 12, 19, 14, 19, 233, 11, 19, 3, 20, 3, 20, 3, 20, 3, 20, 
	3, 20, 3, 20, 5, 20, 241, 10, 20, 3, 21, 3, 21, 3, 22, 3, 22, 3, 22, 3, 
	22, 3, 23, 3, 23, 7, 23, 251, 10, 23, 12, 23, 14, 23, 254, 11, 23, 3, 24, 
	3, 24, 3, 24, 7, 24, 259, 10, 24, 12, 24, 14, 24, 262, 11, 24, 3, 25, 3, 
	25, 3, 25, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 5, 26, 273, 10, 26, 
	3, 26, 3, 26, 3, 26, 3, 26, 7, 26, 279, 10, 26, 12, 26, 14, 26, 282, 11, 
	26, 3, 27, 3, 27, 3, 28, 3, 28, 3, 29, 3, 29, 3, 29, 3, 29, 3, 30, 3, 30, 
	3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 5, 30, 300, 10, 30, 3, 31, 3, 
	31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 5, 31, 310, 10, 31, 3, 31, 
	3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 
	31, 7, 31, 324, 10, 31, 12, 31, 14, 31, 327, 11, 31, 3, 32, 3, 32, 3, 32, 
	3, 33, 3, 33, 3, 34, 3, 34, 3, 34, 5, 34, 337, 10, 34, 3, 34, 3, 34, 3, 
	35, 3, 35, 3, 35, 7, 35, 344, 10, 35, 12, 35, 14, 35, 347, 11, 35, 3, 36, 
	3, 36, 5, 36, 351, 10, 36, 3, 37, 3, 37, 5, 37, 355, 10, 37, 3, 37, 3, 
	37, 5, 37, 359, 10, 37, 3, 38, 3, 38, 3, 38, 3, 38, 3, 39, 5, 39, 366, 
	10, 39, 3, 39, 3, 39, 3, 40, 5, 40, 371, 10, 40, 3, 40, 3, 40, 3, 41, 3, 
	41, 3, 41, 3, 42, 3, 42, 3, 43, 3, 43, 3, 44, 3, 44, 3, 45, 3, 45, 5, 45, 
	386, 10, 45, 3, 45, 3, 45, 3, 45, 5, 45, 391, 10, 45, 7, 45, 393, 10, 45, 
	12, 45, 14, 45, 396, 11, 45, 3, 46, 3, 46, 3, 46, 2, 5, 22, 50, 60, 47, 
	2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 
	40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 
	76, 78, 80, 82, 84, 86, 88, 90, 2, 9, 3, 2, 40, 41, 4, 2, 43, 44, 93, 94, 
	3, 2, 46, 47, 4, 2, 48, 48, 78, 78, 3, 2, 62, 68, 3, 2, 87, 88, 11, 2, 
	3, 3, 7, 7, 9, 11, 15, 24, 26, 29, 31, 35, 38, 52, 54, 57, 61, 68, 2, 411, 
	2, 92, 3, 2, 2, 2, 4, 95, 3, 2, 2, 2, 6, 98, 3, 2, 2, 2, 8, 117, 3, 2, 
	2, 2, 10, 120, 3, 2, 2, 2, 12, 128, 3, 2, 2, 2, 14, 132, 3, 2, 2, 2, 16, 
	135, 3, 2, 2, 2, 18, 138, 3, 2, 2, 2, 20, 143, 3, 2, 2, 2, 22, 173, 3, 
	2, 2, 2, 24, 183, 3, 2, 2, 2, 26, 191, 3, 2, 2, 2, 28, 196, 3, 2, 2, 2, 
	30, 202, 3, 2, 2, 2, 32, 206, 3, 2, 2, 2, 34, 213, 3, 2, 2, 2, 36, 226, 
	3, 2, 2, 2, 38, 240, 3, 2, 2, 2, 40, 242, 3, 2, 2, 2, 42, 244, 3, 2, 2, 
	2, 44, 248, 3, 2, 2, 2, 46, 255, 3, 2, 2, 2, 48, 263, 3, 2, 2, 2, 50, 272, 
	3, 2, 2, 2, 52, 283, 3, 2, 2, 2, 54, 285, 3, 2, 2, 2, 56, 287, 3, 2, 2, 
	2, 58, 299, 3, 2, 2, 2, 60, 309, 3, 2, 2, 2, 62, 328, 3, 2, 2, 2, 64, 331, 
	3, 2, 2, 2, 66, 333, 3, 2, 2, 2, 68, 340, 3, 2, 2, 2, 70, 350, 3, 2, 2, 
	2, 72, 358, 3, 2, 2, 2, 74, 360, 3, 2, 2, 2, 76, 365, 3, 2, 2, 2, 78, 370, 
	3, 2, 2, 2, 80, 374, 3, 2, 2, 2, 82, 377, 3, 2, 2, 2, 84, 379, 3, 2, 2, 
	2, 86, 381, 3, 2, 2, 2, 88, 385, 3, 2, 2, 2, 90, 397, 3, 2, 2, 2, 92, 93, 
	5, 4, 3, 2, 93, 94, 7, 2, 2, 3, 94, 3, 3, 2, 2, 2, 95, 96, 5, 6, 4, 2, 
	96, 5, 3, 2, 2, 2, 97, 99, 7, 36, 2, 2, 98, 97, 3, 2, 2, 2, 98, 99, 3, 
	2, 2, 2, 99, 100, 3, 2, 2, 2, 100, 101, 5, 8, 5, 2, 101, 103, 5, 16, 9, 
	2, 102, 104, 5, 18, 10, 2, 103, 102, 3, 2, 2, 2, 103, 104, 3, 2, 2, 2, 
	104, 106, 3, 2, 2, 2, 105, 107, 5, 34, 18, 2, 106, 105, 3, 2, 2, 2, 106, 
	107, 3, 2, 2, 2, 107, 109, 3, 2, 2, 2, 108, 110, 5, 42, 22, 2, 109, 108, 
	3, 2, 2, 2, 109, 110, 3, 2, 2, 2, 110, 112, 3, 2, 2, 2, 111, 113, 5, 80, 
	41, 2, 112, 111, 3, 2, 2, 2, 112, 113, 3, 2, 2, 2, 113, 115, 3, 2, 2, 2, 
	114, 116, 7, 37, 2, 2, 115, 114, 3, 2, 2, 2, 115, 116, 3, 2, 2, 2, 116, 
	7, 3, 2, 2, 2, 117, 118, 7, 38, 2, 2, 118, 119, 5, 10, 6, 2, 119, 9, 3, 
	2, 2, 2, 120, 125, 5, 12, 7, 2, 121, 122, 7, 80, 2, 2, 122, 124, 5, 12, 
	7, 2, 123, 121, 3, 2, 2, 2, 124, 127, 3, 2, 2, 2, 125, 123, 3, 2, 2, 2, 
	125, 126, 3, 2, 2, 2, 126, 11, 3, 2, 2, 2, 127, 125, 3, 2, 2, 2, 128, 130, 
	5, 60, 31, 2, 129, 131, 5, 14, 8, 2, 130, 129, 3, 2, 2, 2, 130, 131, 3, 
	2, 2, 2, 131, 13, 3, 2, 2, 2, 132, 133, 7, 39, 2, 2, 133, 134, 5, 88, 45, 
	2, 134, 15, 3, 2, 2, 2, 135, 136, 7, 31, 2, 2, 136, 137, 5, 82, 42, 2, 
	137, 17, 3, 2, 2, 2, 138, 139, 7, 32, 2, 2, 139, 140, 5, 20, 11, 2, 140, 
	19, 3, 2, 2, 2, 141, 144, 5, 22, 12, 2, 142, 144, 5, 26, 14, 2, 143, 141, 
	3, 2, 2, 2, 143, 142, 3, 2, 2, 2, 144, 21, 3, 2, 2, 2, 145, 146, 8, 12, 
	1, 2, 146, 147, 7, 85, 2, 2, 147, 148, 5, 22, 12, 2, 148, 149, 7, 86, 2, 
	2, 149, 174, 3, 2, 2, 2, 150, 159, 5, 84, 43, 2, 151, 160, 7, 71, 2, 2, 
	152, 160, 7, 48, 2, 2, 153, 154, 7, 49, 2, 2, 154, 160, 7, 48, 2, 2, 155, 
	160, 7, 78, 2, 2, 156, 160, 7, 79, 2, 2, 157, 160, 7, 72, 2, 2, 158, 160, 
	7, 73, 2, 2, 159, 151, 3, 2, 2, 2, 159, 152, 3, 2, 2, 2, 159, 153, 3, 2, 
	2, 2, 159, 155, 3, 2, 2, 2, 159, 156, 3, 2, 2, 2, 159, 157, 3, 2, 2, 2, 
	159, 158, 3, 2, 2, 2, 160, 161, 3, 2, 2, 2, 161, 162, 5, 86, 44, 2, 162, 
	174, 3, 2, 2, 2, 163, 167, 5, 84, 43, 2, 164, 168, 7, 59, 2, 2, 165, 166, 
	7, 49, 2, 2, 166, 168, 7, 59, 2, 2, 167, 164, 3, 2, 2, 2, 167, 165, 3, 
	2, 2, 2, 168, 169, 3, 2, 2, 2, 169, 170, 7, 85, 2, 2, 170, 171, 5, 24, 
	13, 2, 171, 172, 7, 86, 2, 2, 172, 174, 3, 2, 2, 2, 173, 145, 3, 2, 2, 
	2, 173, 150, 3, 2, 2, 2, 173, 163, 3, 2, 2, 2, 174, 180, 3, 2, 2, 2, 175, 
	176, 12, 3, 2, 2, 176, 177, 9, 2, 2, 2, 177, 179, 5, 22, 12, 4, 178, 175, 
	3, 2, 2, 2, 179, 182, 3, 2, 2, 2, 180, 178, 3, 2, 2, 2, 180, 181, 3, 2, 
	2, 2, 181, 23, 3, 2, 2, 2, 182, 180, 3, 2, 2, 2, 183, 188, 5, 86, 44, 2, 
	184, 185, 7, 80, 2, 2, 185, 187, 5, 86, 44, 2, 186, 184, 3, 2, 2, 2, 187, 
	190, 3, 2, 2, 2, 188, 186, 3, 2, 2, 2, 188, 189, 3, 2, 2, 2, 189, 25, 3, 
	2, 2, 2, 190, 188, 3, 2, 2, 2, 191, 194, 5, 28, 15, 2, 192, 193, 7, 40, 
	2, 2, 193, 195, 5, 28, 15, 2, 194, 192, 3, 2, 2, 2, 194, 195, 3, 2, 2, 
	2, 195, 27, 3, 2, 2, 2, 196, 197, 7, 57, 2, 2, 197, 200, 5, 58, 30, 2, 
	198, 201, 5, 30, 16, 2, 199, 201, 5, 88, 45, 2, 200, 198, 3, 2, 2, 2, 200, 
	199, 3, 2, 2, 2, 201, 29, 3, 2, 2, 2, 202, 204, 5, 32, 17, 2, 203, 205, 
	5, 62, 32, 2, 204, 203, 3, 2, 2, 2, 204, 205, 3, 2, 2, 2, 205, 31, 3, 2, 
	2, 2, 206, 207, 7, 58, 2, 2, 207, 209, 7, 85, 2, 2, 208, 210, 5, 68, 35, 
	2, 209, 208, 3, 2, 2, 2, 209, 210, 3, 2, 2, 2, 210, 211, 3, 2, 2, 2, 211, 
	212, 7, 86, 2, 2, 212, 33, 3, 2, 2, 2, 213, 214, 7, 52, 2, 2, 214, 215, 
	7, 54, 2, 2, 215, 221, 5, 36, 19, 2, 216, 217, 7, 42, 2, 2, 217, 218, 7, 
	85, 2, 2, 218, 219, 5, 40, 21, 2, 219, 220, 7, 86, 2, 2, 220, 222, 3, 2, 
	2, 2, 221, 216, 3, 2, 2, 2, 221, 222, 3, 2, 2, 2, 222, 224, 3, 2, 2, 2, 
	223, 225, 5, 48, 25, 2, 224, 223, 3, 2, 2, 2, 224, 225, 3, 2, 2, 2, 225, 
	35, 3, 2, 2, 2, 226, 231, 5, 38, 20, 2, 227, 228, 7, 80, 2, 2, 228, 230, 
	5, 38, 20, 2, 229, 227, 3, 2, 2, 2, 230, 233, 3, 2, 2, 2, 231, 229, 3, 
	2, 2, 2, 231, 232, 3, 2, 2, 2, 232, 37, 3, 2, 2, 2, 233, 231, 3, 2, 2, 
	2, 234, 241, 5, 88, 45, 2, 235, 236, 7, 57, 2, 2, 236, 237, 7, 85, 2, 2, 
	237, 238, 5, 62, 32, 2, 238, 239, 7, 86, 2, 2, 239, 241, 3, 2, 2, 2, 240, 
	234, 3, 2, 2, 2, 240, 235, 3, 2, 2, 2, 241, 39, 3, 2, 2, 2, 242, 243, 9, 
	3, 2, 2, 243, 41, 3, 2, 2, 2, 244, 245, 7, 45, 2, 2, 245, 246, 7, 54, 2, 
	2, 246, 247, 5, 46, 24, 2, 247, 43, 3, 2, 2, 2, 248, 252, 5, 60, 31, 2, 
	249, 251, 9, 4, 2, 2, 250, 249, 3, 2, 2, 2, 251, 254, 3, 2, 2, 2, 252, 
	250, 3, 2, 2, 2, 252, 253, 3, 2, 2, 2, 253, 45, 3, 2, 2, 2, 254, 252, 3, 
	2, 2, 2, 255, 260, 5, 44, 23, 2, 256, 257, 7, 80, 2, 2, 257, 259, 5, 44, 
	23, 2, 258, 256, 3, 2, 2, 2, 259, 262, 3, 2, 2, 2, 260, 258, 3, 2, 2, 2, 
	260, 261, 3, 2, 2, 2, 261, 47, 3, 2, 2, 2, 262, 260, 3, 2, 2, 2, 263, 264, 
	7, 53, 2, 2, 264, 265, 5, 50, 26, 2, 265, 49, 3, 2, 2, 2, 266, 267, 8, 
	26, 1, 2, 267, 268, 7, 85, 2, 2, 268, 269, 5, 50, 26, 2, 269, 270, 7, 86, 
	2, 2, 270, 273, 3, 2, 2, 2, 271, 273, 5, 54, 28, 2, 272, 266, 3, 2, 2, 
	2, 272, 271, 3, 2, 2, 2, 273, 280, 3, 2, 2, 2, 274, 275, 12, 4, 2, 2, 275, 
	276, 5, 52, 27, 2, 276, 277, 5, 50, 26, 5, 277, 279, 3, 2, 2, 2, 278, 274, 
	3, 2, 2, 2, 279, 282, 3, 2, 2, 2, 280, 278, 3, 2, 2, 2, 280, 281, 3, 2, 
	2, 2, 281, 51, 3, 2, 2, 2, 282, 280, 3, 2, 2, 2, 283, 284, 9, 2, 2, 2, 
	284, 53, 3, 2, 2, 2, 285, 286, 5, 56, 29, 2, 286, 55, 3, 2, 2, 2, 287, 
	288, 5, 60, 31, 2, 288, 289, 5, 58, 30, 2, 289, 290, 5, 60, 31, 2, 290, 
	57, 3, 2, 2, 2, 291, 300, 7, 71, 2, 2, 292, 300, 7, 72, 2, 2, 293, 300, 
	7, 73, 2, 2, 294, 300, 7, 76, 2, 2, 295, 300, 7, 77, 2, 2, 296, 300, 7, 
	74, 2, 2, 297, 300, 7, 75, 2, 2, 298, 300, 9, 5, 2, 2, 299, 291, 3, 2, 
	2, 2, 299, 292, 3, 2, 2, 2, 299, 293, 3, 2, 2, 2, 299, 294, 3, 2, 2, 2, 
	299, 295, 3, 2, 2, 2, 299, 296, 3, 2, 2, 2, 299, 297, 3, 2, 2, 2, 299, 
	298, 3, 2, 2, 2, 300, 59, 3, 2, 2, 2, 301, 302, 8, 31, 1, 2, 302, 303, 
	7, 85, 2, 2, 303, 304, 5, 60, 31, 2, 304, 305, 7, 86, 2, 2, 305, 310, 3, 
	2, 2, 2, 306, 310, 5, 66, 34, 2, 307, 310, 5, 72, 37, 2, 308, 310, 5, 62, 
	32, 2, 309, 301, 3, 2, 2, 2, 309, 306, 3, 2, 2, 2, 309, 307, 3, 2, 2, 2, 
	309, 308, 3, 2, 2, 2, 310, 325, 3, 2, 2, 2, 311, 312, 12, 10, 2, 2, 312, 
	313, 7, 90, 2, 2, 313, 324, 5, 60, 31, 11, 314, 315, 12, 9, 2, 2, 315, 
	316, 7, 89, 2, 2, 316, 324, 5, 60, 31, 10, 317, 318, 12, 8, 2, 2, 318, 
	319, 7, 87, 2, 2, 319, 324, 5, 60, 31, 9, 320, 321, 12, 7, 2, 2, 321, 322, 
	7, 88, 2, 2, 322, 324, 5, 60, 31, 8, 323, 311, 3, 2, 2, 2, 323, 314, 3, 
	2, 2, 2, 323, 317, 3, 2, 2, 2, 323, 320, 3, 2, 2, 2, 324, 327, 3, 2, 2, 
	2, 325, 323, 3, 2, 2, 2, 325, 326, 3, 2, 2, 2, 326, 61, 3, 2, 2, 2, 327, 
	325, 3, 2, 2, 2, 328, 329, 5, 76, 39, 2, 329, 330, 5, 64, 33, 2, 330, 63, 
	3, 2, 2, 2, 331, 332, 9, 6, 2, 2, 332, 65, 3, 2, 2, 2, 333, 334, 5, 88, 
	45, 2, 334, 336, 7, 85, 2, 2, 335, 337, 5, 68, 35, 2, 336, 335, 3, 2, 2, 
	2, 336, 337, 3, 2, 2, 2, 337, 338, 3, 2, 2, 2, 338, 339, 7, 86, 2, 2, 339, 
	67, 3, 2, 2, 2, 340, 345, 5, 70, 36, 2, 341, 342, 7, 80, 2, 2, 342, 344, 
	5, 70, 36, 2, 343, 341, 3, 2, 2, 2, 344, 347, 3, 2, 2, 2, 345, 343, 3, 
	2, 2, 2, 345, 346, 3, 2, 2, 2, 346, 69, 3, 2, 2, 2, 347, 345, 3, 2, 2, 
	2, 348, 351, 5, 60, 31, 2, 349, 351, 5, 22, 12, 2, 350, 348, 3, 2, 2, 2, 
	350, 349, 3, 2, 2, 2, 351, 71, 3, 2, 2, 2, 352, 354, 5, 88, 45, 2, 353, 
	355, 5, 74, 38, 2, 354, 353, 3, 2, 2, 2, 354, 355, 3, 2, 2, 2, 355, 359, 
	3, 2, 2, 2, 356, 359, 5, 78, 40, 2, 357, 359, 5, 76, 39, 2, 358, 352, 3, 
	2, 2, 2, 358, 356, 3, 2, 2, 2, 358, 357, 3, 2, 2, 2, 359, 73, 3, 2, 2, 
	2, 360, 361, 7, 83, 2, 2, 361, 362, 5, 22, 12, 2, 362, 363, 7, 84, 2, 2, 
	363, 75, 3, 2, 2, 2, 364, 366, 9, 7, 2, 2, 365, 364, 3, 2, 2, 2, 365, 366, 
	3, 2, 2, 2, 366, 367, 3, 2, 2, 2, 367, 368, 7, 93, 2, 2, 368, 77, 3, 2, 
	2, 2, 369, 371, 9, 7, 2, 2, 370, 369, 3, 2, 2, 2, 370, 371, 3, 2, 2, 2, 
	371, 372, 3, 2, 2, 2, 372, 373, 7, 94, 2, 2, 373, 79, 3, 2, 2, 2, 374, 
	375, 7, 33, 2, 2, 375, 376, 7, 93, 2, 2, 376, 81, 3, 2, 2, 2, 377, 378, 
	5, 88, 45, 2, 378, 83, 3, 2, 2, 2, 379, 380, 5, 88, 45, 2, 380, 85, 3, 
	2, 2, 2, 381, 382, 5, 88, 45, 2, 382, 87, 3, 2, 2, 2, 383, 386, 7, 92, 
	2, 2, 384, 386, 5, 90, 46, 2, 385, 383, 3, 2, 2, 2, 385, 384, 3, 2, 2, 
	2, 386, 394, 3, 2, 2, 2, 387, 390, 7, 69, 2, 2, 388, 391, 7, 92, 2, 2, 
	389, 391, 5, 90, 46, 2, 390, 388, 3, 2, 2, 2, 390, 389, 3, 2, 2, 2, 391, 
	393, 3, 2, 2, 2, 392, 387, 3, 2, 2, 2, 393, 396, 3, 2, 2, 2, 394, 392, 
	3, 2, 2, 2, 394, 395, 3, 2, 2, 2, 395, 89, 3, 2, 2, 2, 396, 394, 3, 2, 
	2, 2, 397, 398, 9, 8, 2, 2, 398, 91, 3, 2, 2, 2, 42, 98, 103, 106, 109, 
	112, 115, 125, 130, 143, 159, 167, 173, 180, 188, 194, 200, 204, 209, 221, 
	224, 231, 240, 252, 260, 272, 280, 299, 309, 323, 325, 336, 345, 350, 354, 
	358, 365, 370, 385, 390, 394,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "'m'", "", "", "", "'M'", "", "'.'", "':'", 
	"'='", "'<>'", "'!='", "'>'", "'>='", "'<'", "'<='", "'=~'", "'!~'", "','", 
	"'{'", "'}'", "'['", "']'", "'('", "')'", "'+'", "'-'", "'/'", "'*'", "'%'",
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
	"T_REGEXP", "T_NEQREGEXP", "T_COMMA", "T_OPEN_B", "T_CLOSE_B", "T_OPEN_SB", 
	"T_CLOSE_SB", "T_OPEN_P", "T_CLOSE_P", "T_ADD", "T_SUB", "T_DIV", "T_MUL", 
	"T_MOD", "L_ID", "L_INT", "L_DEC", "WS",
}

var ruleNames = []string{
	"statement", "statementList", "queryStmt", "selectExpr", "fields", "field", 
	"alias", "fromClause", "whereClause", "conditionExpr", "tagFilterExpr", 
	"tagValueList", "timeRangeExpr", "timeExpr", "nowExpr", "nowFunc", "groupByClause", 
	"groupByKeys", "groupByKey", "fillOption", "orderByClause", "sortField", 
	"sortFields", "havingClause", "boolExpr", "boolExprLogicalOp", "boolExprAtom", 
	"binaryExpr", "binaryOperator", "fieldExpr", "durationLit", "intervalItem", 
	"exprFunc", "exprFuncParams", "funcParam", "exprAtom", "identFilter", "intNumber", 
	"decNumber", "limitClause", "metricName", "tagKey", "tagValue", "ident", 
	"nonReservedWords",
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
	SQLParserEOF = antlr.TokenEOF
	SQLParserT_CREATE = 1
	SQLParserT_UPDATE = 2
	SQLParserT_SET = 3
	SQLParserT_DROP = 4
	SQLParserT_INTERVAL = 5
	SQLParserT_INTERVAL_NAME = 6
	SQLParserT_SHARD = 7
	SQLParserT_REPLICATION = 8
	SQLParserT_TTL = 9
	SQLParserT_META_TTL = 10
	SQLParserT_PAST_TTL = 11
	SQLParserT_FUTURE_TTL = 12
	SQLParserT_KILL = 13
	SQLParserT_ON = 14
	SQLParserT_SHOW = 15
	SQLParserT_DATASBAE = 16
	SQLParserT_DATASBAES = 17
	SQLParserT_NODE = 18
	SQLParserT_MEASUREMENTS = 19
	SQLParserT_MEASUREMENT = 20
	SQLParserT_FIELD = 21
	SQLParserT_TAG = 22
	SQLParserT_INFO = 23
	SQLParserT_KEYS = 24
	SQLParserT_KEY = 25
	SQLParserT_WITH = 26
	SQLParserT_VALUES = 27
	SQLParserT_VALUE = 28
	SQLParserT_FROM = 29
	SQLParserT_WHERE = 30
	SQLParserT_LIMIT = 31
	SQLParserT_QUERIES = 32
	SQLParserT_QUERY = 33
	SQLParserT_EXPLAIN = 34
	SQLParserT_WITH_VALUE = 35
	SQLParserT_SELECT = 36
	SQLParserT_AS = 37
	SQLParserT_AND = 38
	SQLParserT_OR = 39
	SQLParserT_FILL = 40
	SQLParserT_NULL = 41
	SQLParserT_PREVIOUS = 42
	SQLParserT_ORDER = 43
	SQLParserT_ASC = 44
	SQLParserT_DESC = 45
	SQLParserT_LIKE = 46
	SQLParserT_NOT = 47
	SQLParserT_BETWEEN = 48
	SQLParserT_IS = 49
	SQLParserT_GROUP = 50
	SQLParserT_HAVING = 51
	SQLParserT_BY = 52
	SQLParserT_FOR = 53
	SQLParserT_STATS = 54
	SQLParserT_TIME = 55
	SQLParserT_NOW = 56
	SQLParserT_IN = 57
	SQLParserT_LOG = 58
	SQLParserT_PROFILE = 59
	SQLParserT_SECOND = 60
	SQLParserT_MINUTE = 61
	SQLParserT_HOUR = 62
	SQLParserT_DAY = 63
	SQLParserT_WEEK = 64
	SQLParserT_MONTH = 65
	SQLParserT_YEAR = 66
	SQLParserT_DOT = 67
	SQLParserT_COLON = 68
	SQLParserT_EQUAL = 69
	SQLParserT_NOTEQUAL = 70
	SQLParserT_NOTEQUAL2 = 71
	SQLParserT_GREATER = 72
	SQLParserT_GREATEREQUAL = 73
	SQLParserT_LESS = 74
	SQLParserT_LESSEQUAL = 75
	SQLParserT_REGEXP = 76
	SQLParserT_NEQREGEXP = 77
	SQLParserT_COMMA = 78
	SQLParserT_OPEN_B = 79
	SQLParserT_CLOSE_B = 80
	SQLParserT_OPEN_SB = 81
	SQLParserT_CLOSE_SB = 82
	SQLParserT_OPEN_P = 83
	SQLParserT_CLOSE_P = 84
	SQLParserT_ADD = 85
	SQLParserT_SUB = 86
	SQLParserT_DIV = 87
	SQLParserT_MUL = 88
	SQLParserT_MOD = 89
	SQLParserL_ID = 90
	SQLParserL_INT = 91
	SQLParserL_DEC = 92
	SQLParserWS = 93
)

// SQLParser rules.
const (
	SQLParserRULE_statement = 0
	SQLParserRULE_statementList = 1
	SQLParserRULE_queryStmt = 2
	SQLParserRULE_selectExpr = 3
	SQLParserRULE_fields = 4
	SQLParserRULE_field = 5
	SQLParserRULE_alias = 6
	SQLParserRULE_fromClause = 7
	SQLParserRULE_whereClause = 8
	SQLParserRULE_conditionExpr = 9
	SQLParserRULE_tagFilterExpr = 10
	SQLParserRULE_tagValueList = 11
	SQLParserRULE_timeRangeExpr = 12
	SQLParserRULE_timeExpr = 13
	SQLParserRULE_nowExpr = 14
	SQLParserRULE_nowFunc = 15
	SQLParserRULE_groupByClause = 16
	SQLParserRULE_groupByKeys = 17
	SQLParserRULE_groupByKey = 18
	SQLParserRULE_fillOption = 19
	SQLParserRULE_orderByClause = 20
	SQLParserRULE_sortField = 21
	SQLParserRULE_sortFields = 22
	SQLParserRULE_havingClause = 23
	SQLParserRULE_boolExpr = 24
	SQLParserRULE_boolExprLogicalOp = 25
	SQLParserRULE_boolExprAtom = 26
	SQLParserRULE_binaryExpr = 27
	SQLParserRULE_binaryOperator = 28
	SQLParserRULE_fieldExpr = 29
	SQLParserRULE_durationLit = 30
	SQLParserRULE_intervalItem = 31
	SQLParserRULE_exprFunc = 32
	SQLParserRULE_exprFuncParams = 33
	SQLParserRULE_funcParam = 34
	SQLParserRULE_exprAtom = 35
	SQLParserRULE_identFilter = 36
	SQLParserRULE_intNumber = 37
	SQLParserRULE_decNumber = 38
	SQLParserRULE_limitClause = 39
	SQLParserRULE_metricName = 40
	SQLParserRULE_tagKey = 41
	SQLParserRULE_tagValue = 42
	SQLParserRULE_ident = 43
	SQLParserRULE_nonReservedWords = 44
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
		p.SetState(90)
		p.StatementList()
	}
	{
		p.SetState(91)
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

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(93)
		p.QueryStmt()
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

func (s *QueryStmtContext) SelectExpr() ISelectExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISelectExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISelectExprContext)
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




func (p *SQLParser) QueryStmt() (localctx IQueryStmtContext) {
	localctx = NewQueryStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SQLParserRULE_queryStmt)
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
	p.SetState(96)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_EXPLAIN {
		{
			p.SetState(95)
			p.Match(SQLParserT_EXPLAIN)
		}

	}
	{
		p.SetState(98)
		p.SelectExpr()
	}
	{
		p.SetState(99)
		p.FromClause()
	}
	p.SetState(101)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(100)
			p.WhereClause()
		}

	}
	p.SetState(104)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_GROUP {
		{
			p.SetState(103)
			p.GroupByClause()
		}

	}
	p.SetState(107)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ORDER {
		{
			p.SetState(106)
			p.OrderByClause()
		}

	}
	p.SetState(110)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(109)
			p.LimitClause()
		}

	}
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WITH_VALUE {
		{
			p.SetState(112)
			p.Match(SQLParserT_WITH_VALUE)
		}

	}



	return localctx
}


// ISelectExprContext is an interface to support dynamic dispatch.
type ISelectExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSelectExprContext differentiates from other interfaces.
	IsSelectExprContext()
}

type SelectExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectExprContext() *SelectExprContext {
	var p = new(SelectExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_selectExpr
	return p
}

func (*SelectExprContext) IsSelectExprContext() {}

func NewSelectExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectExprContext {
	var p = new(SelectExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_selectExpr

	return p
}

func (s *SelectExprContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectExprContext) T_SELECT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SELECT, 0)
}

func (s *SelectExprContext) Fields() IFieldsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldsContext)
}

func (s *SelectExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SelectExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSelectExpr(s)
	}
}

func (s *SelectExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSelectExpr(s)
	}
}




func (p *SQLParser) SelectExpr() (localctx ISelectExprContext) {
	localctx = NewSelectExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, SQLParserRULE_selectExpr)

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
		p.SetState(115)
		p.Match(SQLParserT_SELECT)
	}
	{
		p.SetState(116)
		p.Fields()
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




func (p *SQLParser) Fields() (localctx IFieldsContext) {
	localctx = NewFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SQLParserRULE_fields)
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
		p.SetState(118)
		p.Field()
	}
	p.SetState(123)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(119)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(120)
			p.Field()
		}


		p.SetState(125)
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

func (s *FieldContext) FieldExpr() IFieldExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldExprContext)
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




func (p *SQLParser) Field() (localctx IFieldContext) {
	localctx = NewFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SQLParserRULE_field)
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
		p.SetState(126)
		p.fieldExpr(0)
	}
	p.SetState(128)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_AS {
		{
			p.SetState(127)
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




func (p *SQLParser) Alias() (localctx IAliasContext) {
	localctx = NewAliasContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SQLParserRULE_alias)

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
		p.SetState(130)
		p.Match(SQLParserT_AS)
	}
	{
		p.SetState(131)
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




func (p *SQLParser) FromClause() (localctx IFromClauseContext) {
	localctx = NewFromClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SQLParserRULE_fromClause)

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
		p.SetState(133)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(134)
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

func (s *WhereClauseContext) ConditionExpr() IConditionExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IConditionExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IConditionExprContext)
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




func (p *SQLParser) WhereClause() (localctx IWhereClauseContext) {
	localctx = NewWhereClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SQLParserRULE_whereClause)

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
		p.SetState(136)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(137)
		p.ConditionExpr()
	}



	return localctx
}


// IConditionExprContext is an interface to support dynamic dispatch.
type IConditionExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsConditionExprContext differentiates from other interfaces.
	IsConditionExprContext()
}

type ConditionExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConditionExprContext() *ConditionExprContext {
	var p = new(ConditionExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_conditionExpr
	return p
}

func (*ConditionExprContext) IsConditionExprContext() {}

func NewConditionExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionExprContext {
	var p = new(ConditionExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_conditionExpr

	return p
}

func (s *ConditionExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionExprContext) TagFilterExpr() ITagFilterExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagFilterExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagFilterExprContext)
}

func (s *ConditionExprContext) TimeRangeExpr() ITimeRangeExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITimeRangeExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITimeRangeExprContext)
}

func (s *ConditionExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConditionExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ConditionExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterConditionExpr(s)
	}
}

func (s *ConditionExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitConditionExpr(s)
	}
}




func (p *SQLParser) ConditionExpr() (localctx IConditionExprContext) {
	localctx = NewConditionExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SQLParserRULE_conditionExpr)

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

	p.SetState(141)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(139)
			p.tagFilterExpr(0)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(140)
			p.TimeRangeExpr()
		}

	}


	return localctx
}


// ITagFilterExprContext is an interface to support dynamic dispatch.
type ITagFilterExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagFilterExprContext differentiates from other interfaces.
	IsTagFilterExprContext()
}

type TagFilterExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagFilterExprContext() *TagFilterExprContext {
	var p = new(TagFilterExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagFilterExpr
	return p
}

func (*TagFilterExprContext) IsTagFilterExprContext() {}

func NewTagFilterExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagFilterExprContext {
	var p = new(TagFilterExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagFilterExpr

	return p
}

func (s *TagFilterExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TagFilterExprContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *TagFilterExprContext) AllTagFilterExpr() []ITagFilterExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITagFilterExprContext)(nil)).Elem())
	var tst = make([]ITagFilterExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITagFilterExprContext)
		}
	}

	return tst
}

func (s *TagFilterExprContext) TagFilterExpr(i int) ITagFilterExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagFilterExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITagFilterExprContext)
}

func (s *TagFilterExprContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *TagFilterExprContext) TagKey() ITagKeyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagKeyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagKeyContext)
}

func (s *TagFilterExprContext) TagValue() ITagValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagValueContext)
}

func (s *TagFilterExprContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *TagFilterExprContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *TagFilterExprContext) T_NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOT, 0)
}

func (s *TagFilterExprContext) T_REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REGEXP, 0)
}

func (s *TagFilterExprContext) T_NEQREGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NEQREGEXP, 0)
}

func (s *TagFilterExprContext) T_NOTEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL, 0)
}

func (s *TagFilterExprContext) T_NOTEQUAL2() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL2, 0)
}

func (s *TagFilterExprContext) TagValueList() ITagValueListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagValueListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagValueListContext)
}

func (s *TagFilterExprContext) T_IN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IN, 0)
}

func (s *TagFilterExprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *TagFilterExprContext) T_OR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OR, 0)
}

func (s *TagFilterExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagFilterExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TagFilterExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagFilterExpr(s)
	}
}

func (s *TagFilterExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagFilterExpr(s)
	}
}





func (p *SQLParser) TagFilterExpr() (localctx ITagFilterExprContext) {
	return p.tagFilterExpr(0)
}

func (p *SQLParser) tagFilterExpr(_p int) (localctx ITagFilterExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewTagFilterExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ITagFilterExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 20
	p.EnterRecursionRule(localctx, 20, SQLParserRULE_tagFilterExpr, _p)
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
	p.SetState(171)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(144)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(145)
			p.tagFilterExpr(0)
		}
		{
			p.SetState(146)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(148)
			p.TagKey()
		}
		p.SetState(157)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_EQUAL:
			{
				p.SetState(149)
				p.Match(SQLParserT_EQUAL)
			}


		case SQLParserT_LIKE:
			{
				p.SetState(150)
				p.Match(SQLParserT_LIKE)
			}


		case SQLParserT_NOT:
			{
				p.SetState(151)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(152)
				p.Match(SQLParserT_LIKE)
			}


		case SQLParserT_REGEXP:
			{
				p.SetState(153)
				p.Match(SQLParserT_REGEXP)
			}


		case SQLParserT_NEQREGEXP:
			{
				p.SetState(154)
				p.Match(SQLParserT_NEQREGEXP)
			}


		case SQLParserT_NOTEQUAL:
			{
				p.SetState(155)
				p.Match(SQLParserT_NOTEQUAL)
			}


		case SQLParserT_NOTEQUAL2:
			{
				p.SetState(156)
				p.Match(SQLParserT_NOTEQUAL2)
			}



		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(159)
			p.TagValue()
		}


	case 3:
		{
			p.SetState(161)
			p.TagKey()
		}
		p.SetState(165)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_IN:
			{
				p.SetState(162)
				p.Match(SQLParserT_IN)
			}


		case SQLParserT_NOT:
			{
				p.SetState(163)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(164)
				p.Match(SQLParserT_IN)
			}



		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(167)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(168)
			p.TagValueList()
		}
		{
			p.SetState(169)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(178)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewTagFilterExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_tagFilterExpr)
			p.SetState(173)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(174)
				_la = p.GetTokenStream().LA(1)

				if !(_la == SQLParserT_AND || _la == SQLParserT_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(175)
				p.tagFilterExpr(2)
			}


		}
		p.SetState(180)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext())
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




func (p *SQLParser) TagValueList() (localctx ITagValueListContext) {
	localctx = NewTagValueListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SQLParserRULE_tagValueList)
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
		p.SetState(181)
		p.TagValue()
	}
	p.SetState(186)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(182)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(183)
			p.TagValue()
		}


		p.SetState(188)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// ITimeRangeExprContext is an interface to support dynamic dispatch.
type ITimeRangeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTimeRangeExprContext differentiates from other interfaces.
	IsTimeRangeExprContext()
}

type TimeRangeExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeRangeExprContext() *TimeRangeExprContext {
	var p = new(TimeRangeExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_timeRangeExpr
	return p
}

func (*TimeRangeExprContext) IsTimeRangeExprContext() {}

func NewTimeRangeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeRangeExprContext {
	var p = new(TimeRangeExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_timeRangeExpr

	return p
}

func (s *TimeRangeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeRangeExprContext) AllTimeExpr() []ITimeExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITimeExprContext)(nil)).Elem())
	var tst = make([]ITimeExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITimeExprContext)
		}
	}

	return tst
}

func (s *TimeRangeExprContext) TimeExpr(i int) ITimeExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITimeExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITimeExprContext)
}

func (s *TimeRangeExprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *TimeRangeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeRangeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TimeRangeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTimeRangeExpr(s)
	}
}

func (s *TimeRangeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTimeRangeExpr(s)
	}
}




func (p *SQLParser) TimeRangeExpr() (localctx ITimeRangeExprContext) {
	localctx = NewTimeRangeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, SQLParserRULE_timeRangeExpr)
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
		p.SetState(189)
		p.TimeExpr()
	}
	p.SetState(192)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_AND {
		{
			p.SetState(190)
			p.Match(SQLParserT_AND)
		}
		{
			p.SetState(191)
			p.TimeExpr()
		}

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

func (s *TimeExprContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *TimeExprContext) BinaryOperator() IBinaryOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBinaryOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBinaryOperatorContext)
}

func (s *TimeExprContext) NowExpr() INowExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INowExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INowExprContext)
}

func (s *TimeExprContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
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




func (p *SQLParser) TimeExpr() (localctx ITimeExprContext) {
	localctx = NewTimeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, SQLParserRULE_timeExpr)

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
		p.SetState(194)
		p.Match(SQLParserT_TIME)
	}
	{
		p.SetState(195)
		p.BinaryOperator()
	}
	p.SetState(198)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_NOW:
		{
			p.SetState(196)
			p.NowExpr()
		}


	case SQLParserT_CREATE, SQLParserT_INTERVAL, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NODE, SQLParserT_MEASUREMENTS, SQLParserT_MEASUREMENT, SQLParserT_FIELD, SQLParserT_TAG, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_PROFILE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR, SQLParserL_ID:
		{
			p.SetState(197)
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




func (p *SQLParser) NowExpr() (localctx INowExprContext) {
	localctx = NewNowExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, SQLParserRULE_nowExpr)
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
		p.SetState(200)
		p.NowFunc()
	}
	p.SetState(202)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if ((((_la - 85)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 85))) & ((1 << (SQLParserT_ADD - 85)) | (1 << (SQLParserT_SUB - 85)) | (1 << (SQLParserL_INT - 85)))) != 0) {
		{
			p.SetState(201)
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




func (p *SQLParser) NowFunc() (localctx INowFuncContext) {
	localctx = NewNowFuncContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, SQLParserRULE_nowFunc)
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
		p.SetState(204)
		p.Match(SQLParserT_NOW)
	}
	{
		p.SetState(205)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(207)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_DATASBAE) | (1 << SQLParserT_DATASBAES) | (1 << SQLParserT_NODE) | (1 << SQLParserT_MEASUREMENTS) | (1 << SQLParserT_MEASUREMENT) | (1 << SQLParserT_FIELD) | (1 << SQLParserT_TAG) | (1 << SQLParserT_KEYS) | (1 << SQLParserT_KEY) | (1 << SQLParserT_WITH) | (1 << SQLParserT_VALUES) | (1 << SQLParserT_FROM) | (1 << SQLParserT_WHERE) | (1 << SQLParserT_LIMIT))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)) | (1 << (SQLParserT_ASC - 32)) | (1 << (SQLParserT_DESC - 32)) | (1 << (SQLParserT_LIKE - 32)) | (1 << (SQLParserT_NOT - 32)) | (1 << (SQLParserT_BETWEEN - 32)) | (1 << (SQLParserT_IS - 32)) | (1 << (SQLParserT_GROUP - 32)) | (1 << (SQLParserT_BY - 32)) | (1 << (SQLParserT_FOR - 32)) | (1 << (SQLParserT_STATS - 32)) | (1 << (SQLParserT_TIME - 32)) | (1 << (SQLParserT_PROFILE - 32)) | (1 << (SQLParserT_SECOND - 32)) | (1 << (SQLParserT_MINUTE - 32)) | (1 << (SQLParserT_HOUR - 32)) | (1 << (SQLParserT_DAY - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)) | (1 << (SQLParserT_OPEN_P - 64)) | (1 << (SQLParserT_ADD - 64)) | (1 << (SQLParserT_SUB - 64)) | (1 << (SQLParserL_ID - 64)) | (1 << (SQLParserL_INT - 64)) | (1 << (SQLParserL_DEC - 64)))) != 0) {
		{
			p.SetState(206)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(209)
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

func (s *GroupByClauseContext) GroupByKeys() IGroupByKeysContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IGroupByKeysContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IGroupByKeysContext)
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




func (p *SQLParser) GroupByClause() (localctx IGroupByClauseContext) {
	localctx = NewGroupByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, SQLParserRULE_groupByClause)
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
		p.SetState(211)
		p.Match(SQLParserT_GROUP)
	}
	{
		p.SetState(212)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(213)
		p.GroupByKeys()
	}
	p.SetState(219)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_FILL {
		{
			p.SetState(214)
			p.Match(SQLParserT_FILL)
		}
		{
			p.SetState(215)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(216)
			p.FillOption()
		}
		{
			p.SetState(217)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.SetState(222)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_HAVING {
		{
			p.SetState(221)
			p.HavingClause()
		}

	}



	return localctx
}


// IGroupByKeysContext is an interface to support dynamic dispatch.
type IGroupByKeysContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGroupByKeysContext differentiates from other interfaces.
	IsGroupByKeysContext()
}

type GroupByKeysContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByKeysContext() *GroupByKeysContext {
	var p = new(GroupByKeysContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_groupByKeys
	return p
}

func (*GroupByKeysContext) IsGroupByKeysContext() {}

func NewGroupByKeysContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByKeysContext {
	var p = new(GroupByKeysContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_groupByKeys

	return p
}

func (s *GroupByKeysContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByKeysContext) AllGroupByKey() []IGroupByKeyContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IGroupByKeyContext)(nil)).Elem())
	var tst = make([]IGroupByKeyContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IGroupByKeyContext)
		}
	}

	return tst
}

func (s *GroupByKeysContext) GroupByKey(i int) IGroupByKeyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IGroupByKeyContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IGroupByKeyContext)
}

func (s *GroupByKeysContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *GroupByKeysContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *GroupByKeysContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByKeysContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *GroupByKeysContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterGroupByKeys(s)
	}
}

func (s *GroupByKeysContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitGroupByKeys(s)
	}
}




func (p *SQLParser) GroupByKeys() (localctx IGroupByKeysContext) {
	localctx = NewGroupByKeysContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, SQLParserRULE_groupByKeys)
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
		p.SetState(224)
		p.GroupByKey()
	}
	p.SetState(229)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(225)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(226)
			p.GroupByKey()
		}


		p.SetState(231)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// IGroupByKeyContext is an interface to support dynamic dispatch.
type IGroupByKeyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGroupByKeyContext differentiates from other interfaces.
	IsGroupByKeyContext()
}

type GroupByKeyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByKeyContext() *GroupByKeyContext {
	var p = new(GroupByKeyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_groupByKey
	return p
}

func (*GroupByKeyContext) IsGroupByKeyContext() {}

func NewGroupByKeyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByKeyContext {
	var p = new(GroupByKeyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_groupByKey

	return p
}

func (s *GroupByKeyContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByKeyContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *GroupByKeyContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *GroupByKeyContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *GroupByKeyContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *GroupByKeyContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *GroupByKeyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByKeyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *GroupByKeyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterGroupByKey(s)
	}
}

func (s *GroupByKeyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitGroupByKey(s)
	}
}




func (p *SQLParser) GroupByKey() (localctx IGroupByKeyContext) {
	localctx = NewGroupByKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, SQLParserRULE_groupByKey)

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

	p.SetState(238)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(232)
			p.Ident()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(233)
			p.Match(SQLParserT_TIME)
		}
		{
			p.SetState(234)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(235)
			p.DurationLit()
		}
		{
			p.SetState(236)
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




func (p *SQLParser) FillOption() (localctx IFillOptionContext) {
	localctx = NewFillOptionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, SQLParserRULE_fillOption)
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
		p.SetState(240)
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




func (p *SQLParser) OrderByClause() (localctx IOrderByClauseContext) {
	localctx = NewOrderByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, SQLParserRULE_orderByClause)

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
		p.Match(SQLParserT_ORDER)
	}
	{
		p.SetState(243)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(244)
		p.SortFields()
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

func (s *SortFieldContext) FieldExpr() IFieldExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldExprContext)
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




func (p *SQLParser) SortField() (localctx ISortFieldContext) {
	localctx = NewSortFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, SQLParserRULE_sortField)
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
		p.SetState(246)
		p.fieldExpr(0)
	}
	p.SetState(250)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_ASC || _la == SQLParserT_DESC {
		{
			p.SetState(247)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_ASC || _la == SQLParserT_DESC) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}


		p.SetState(252)
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




func (p *SQLParser) SortFields() (localctx ISortFieldsContext) {
	localctx = NewSortFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, SQLParserRULE_sortFields)
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
		p.SetState(253)
		p.SortField()
	}
	p.SetState(258)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(254)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(255)
			p.SortField()
		}


		p.SetState(260)
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




func (p *SQLParser) HavingClause() (localctx IHavingClauseContext) {
	localctx = NewHavingClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, SQLParserRULE_havingClause)

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
		p.SetState(261)
		p.Match(SQLParserT_HAVING)
	}
	{
		p.SetState(262)
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





func (p *SQLParser) BoolExpr() (localctx IBoolExprContext) {
	return p.boolExpr(0)
}

func (p *SQLParser) boolExpr(_p int) (localctx IBoolExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewBoolExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IBoolExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 48
	p.EnterRecursionRule(localctx, 48, SQLParserRULE_boolExpr, _p)

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
	p.SetState(270)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(265)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(266)
			p.boolExpr(0)
		}
		{
			p.SetState(267)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(269)
			p.BoolExprAtom()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(278)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewBoolExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_boolExpr)
			p.SetState(272)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(273)
				p.BoolExprLogicalOp()
			}
			{
				p.SetState(274)
				p.boolExpr(3)
			}


		}
		p.SetState(280)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())
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




func (p *SQLParser) BoolExprLogicalOp() (localctx IBoolExprLogicalOpContext) {
	localctx = NewBoolExprLogicalOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, SQLParserRULE_boolExprLogicalOp)
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

func (s *BoolExprAtomContext) BinaryExpr() IBinaryExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBinaryExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBinaryExprContext)
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




func (p *SQLParser) BoolExprAtom() (localctx IBoolExprAtomContext) {
	localctx = NewBoolExprAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, SQLParserRULE_boolExprAtom)

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
		p.SetState(283)
		p.BinaryExpr()
	}



	return localctx
}


// IBinaryExprContext is an interface to support dynamic dispatch.
type IBinaryExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBinaryExprContext differentiates from other interfaces.
	IsBinaryExprContext()
}

type BinaryExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBinaryExprContext() *BinaryExprContext {
	var p = new(BinaryExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_binaryExpr
	return p
}

func (*BinaryExprContext) IsBinaryExprContext() {}

func NewBinaryExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinaryExprContext {
	var p = new(BinaryExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_binaryExpr

	return p
}

func (s *BinaryExprContext) GetParser() antlr.Parser { return s.parser }

func (s *BinaryExprContext) AllFieldExpr() []IFieldExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFieldExprContext)(nil)).Elem())
	var tst = make([]IFieldExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFieldExprContext)
		}
	}

	return tst
}

func (s *BinaryExprContext) FieldExpr(i int) IFieldExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFieldExprContext)
}

func (s *BinaryExprContext) BinaryOperator() IBinaryOperatorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBinaryOperatorContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBinaryOperatorContext)
}

func (s *BinaryExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *BinaryExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBinaryExpr(s)
	}
}

func (s *BinaryExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBinaryExpr(s)
	}
}




func (p *SQLParser) BinaryExpr() (localctx IBinaryExprContext) {
	localctx = NewBinaryExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, SQLParserRULE_binaryExpr)

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
		p.SetState(285)
		p.fieldExpr(0)
	}
	{
		p.SetState(286)
		p.BinaryOperator()
	}
	{
		p.SetState(287)
		p.fieldExpr(0)
	}



	return localctx
}


// IBinaryOperatorContext is an interface to support dynamic dispatch.
type IBinaryOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBinaryOperatorContext differentiates from other interfaces.
	IsBinaryOperatorContext()
}

type BinaryOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBinaryOperatorContext() *BinaryOperatorContext {
	var p = new(BinaryOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_binaryOperator
	return p
}

func (*BinaryOperatorContext) IsBinaryOperatorContext() {}

func NewBinaryOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinaryOperatorContext {
	var p = new(BinaryOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_binaryOperator

	return p
}

func (s *BinaryOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *BinaryOperatorContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *BinaryOperatorContext) T_NOTEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL, 0)
}

func (s *BinaryOperatorContext) T_NOTEQUAL2() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL2, 0)
}

func (s *BinaryOperatorContext) T_LESS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LESS, 0)
}

func (s *BinaryOperatorContext) T_LESSEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LESSEQUAL, 0)
}

func (s *BinaryOperatorContext) T_GREATER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GREATER, 0)
}

func (s *BinaryOperatorContext) T_GREATEREQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GREATEREQUAL, 0)
}

func (s *BinaryOperatorContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *BinaryOperatorContext) T_REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REGEXP, 0)
}

func (s *BinaryOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *BinaryOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBinaryOperator(s)
	}
}

func (s *BinaryOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBinaryOperator(s)
	}
}




func (p *SQLParser) BinaryOperator() (localctx IBinaryOperatorContext) {
	localctx = NewBinaryOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, SQLParserRULE_binaryOperator)
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

	p.SetState(297)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_EQUAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(289)
			p.Match(SQLParserT_EQUAL)
		}


	case SQLParserT_NOTEQUAL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(290)
			p.Match(SQLParserT_NOTEQUAL)
		}


	case SQLParserT_NOTEQUAL2:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(291)
			p.Match(SQLParserT_NOTEQUAL2)
		}


	case SQLParserT_LESS:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(292)
			p.Match(SQLParserT_LESS)
		}


	case SQLParserT_LESSEQUAL:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(293)
			p.Match(SQLParserT_LESSEQUAL)
		}


	case SQLParserT_GREATER:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(294)
			p.Match(SQLParserT_GREATER)
		}


	case SQLParserT_GREATEREQUAL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(295)
			p.Match(SQLParserT_GREATEREQUAL)
		}


	case SQLParserT_LIKE, SQLParserT_REGEXP:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(296)
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


// IFieldExprContext is an interface to support dynamic dispatch.
type IFieldExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldExprContext differentiates from other interfaces.
	IsFieldExprContext()
}

type FieldExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldExprContext() *FieldExprContext {
	var p = new(FieldExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_fieldExpr
	return p
}

func (*FieldExprContext) IsFieldExprContext() {}

func NewFieldExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldExprContext {
	var p = new(FieldExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_fieldExpr

	return p
}

func (s *FieldExprContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldExprContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *FieldExprContext) AllFieldExpr() []IFieldExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFieldExprContext)(nil)).Elem())
	var tst = make([]IFieldExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFieldExprContext)
		}
	}

	return tst
}

func (s *FieldExprContext) FieldExpr(i int) IFieldExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFieldExprContext)
}

func (s *FieldExprContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *FieldExprContext) ExprFunc() IExprFuncContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprFuncContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprFuncContext)
}

func (s *FieldExprContext) ExprAtom() IExprAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprAtomContext)
}

func (s *FieldExprContext) DurationLit() IDurationLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDurationLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *FieldExprContext) T_MUL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MUL, 0)
}

func (s *FieldExprContext) T_DIV() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DIV, 0)
}

func (s *FieldExprContext) T_ADD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ADD, 0)
}

func (s *FieldExprContext) T_SUB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUB, 0)
}

func (s *FieldExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FieldExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFieldExpr(s)
	}
}

func (s *FieldExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFieldExpr(s)
	}
}





func (p *SQLParser) FieldExpr() (localctx IFieldExprContext) {
	return p.fieldExpr(0)
}

func (p *SQLParser) fieldExpr(_p int) (localctx IFieldExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewFieldExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IFieldExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 58
	p.EnterRecursionRule(localctx, 58, SQLParserRULE_fieldExpr, _p)

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
	p.SetState(307)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(300)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(301)
			p.fieldExpr(0)
		}
		{
			p.SetState(302)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(304)
			p.ExprFunc()
		}


	case 3:
		{
			p.SetState(305)
			p.ExprAtom()
		}


	case 4:
		{
			p.SetState(306)
			p.DurationLit()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(323)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(321)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext()) {
			case 1:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(309)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(310)
					p.Match(SQLParserT_MUL)
				}
				{
					p.SetState(311)
					p.fieldExpr(9)
				}


			case 2:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(312)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(313)
					p.Match(SQLParserT_DIV)
				}
				{
					p.SetState(314)
					p.fieldExpr(8)
				}


			case 3:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(315)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				}
				{
					p.SetState(316)
					p.Match(SQLParserT_ADD)
				}
				{
					p.SetState(317)
					p.fieldExpr(7)
				}


			case 4:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(318)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				{
					p.SetState(319)
					p.Match(SQLParserT_SUB)
				}
				{
					p.SetState(320)
					p.fieldExpr(6)
				}

			}

		}
		p.SetState(325)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())
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




func (p *SQLParser) DurationLit() (localctx IDurationLitContext) {
	localctx = NewDurationLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, SQLParserRULE_durationLit)

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
		p.SetState(326)
		p.IntNumber()
	}
	{
		p.SetState(327)
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




func (p *SQLParser) IntervalItem() (localctx IIntervalItemContext) {
	localctx = NewIntervalItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, SQLParserRULE_intervalItem)
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
		p.SetState(329)
		_la = p.GetTokenStream().LA(1)

		if !(((((_la - 60)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 60))) & ((1 << (SQLParserT_SECOND - 60)) | (1 << (SQLParserT_MINUTE - 60)) | (1 << (SQLParserT_HOUR - 60)) | (1 << (SQLParserT_DAY - 60)) | (1 << (SQLParserT_WEEK - 60)) | (1 << (SQLParserT_MONTH - 60)) | (1 << (SQLParserT_YEAR - 60)))) != 0)) {
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




func (p *SQLParser) ExprFunc() (localctx IExprFuncContext) {
	localctx = NewExprFuncContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, SQLParserRULE_exprFunc)
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
		p.SetState(331)
		p.Ident()
	}
	{
		p.SetState(332)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(334)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_DATASBAE) | (1 << SQLParserT_DATASBAES) | (1 << SQLParserT_NODE) | (1 << SQLParserT_MEASUREMENTS) | (1 << SQLParserT_MEASUREMENT) | (1 << SQLParserT_FIELD) | (1 << SQLParserT_TAG) | (1 << SQLParserT_KEYS) | (1 << SQLParserT_KEY) | (1 << SQLParserT_WITH) | (1 << SQLParserT_VALUES) | (1 << SQLParserT_FROM) | (1 << SQLParserT_WHERE) | (1 << SQLParserT_LIMIT))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)) | (1 << (SQLParserT_ASC - 32)) | (1 << (SQLParserT_DESC - 32)) | (1 << (SQLParserT_LIKE - 32)) | (1 << (SQLParserT_NOT - 32)) | (1 << (SQLParserT_BETWEEN - 32)) | (1 << (SQLParserT_IS - 32)) | (1 << (SQLParserT_GROUP - 32)) | (1 << (SQLParserT_BY - 32)) | (1 << (SQLParserT_FOR - 32)) | (1 << (SQLParserT_STATS - 32)) | (1 << (SQLParserT_TIME - 32)) | (1 << (SQLParserT_PROFILE - 32)) | (1 << (SQLParserT_SECOND - 32)) | (1 << (SQLParserT_MINUTE - 32)) | (1 << (SQLParserT_HOUR - 32)) | (1 << (SQLParserT_DAY - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)) | (1 << (SQLParserT_OPEN_P - 64)) | (1 << (SQLParserT_ADD - 64)) | (1 << (SQLParserT_SUB - 64)) | (1 << (SQLParserL_ID - 64)) | (1 << (SQLParserL_INT - 64)) | (1 << (SQLParserL_DEC - 64)))) != 0) {
		{
			p.SetState(333)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(336)
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




func (p *SQLParser) ExprFuncParams() (localctx IExprFuncParamsContext) {
	localctx = NewExprFuncParamsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, SQLParserRULE_exprFuncParams)
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
		p.SetState(338)
		p.FuncParam()
	}
	p.SetState(343)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(339)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(340)
			p.FuncParam()
		}


		p.SetState(345)
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

func (s *FuncParamContext) FieldExpr() IFieldExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldExprContext)
}

func (s *FuncParamContext) TagFilterExpr() ITagFilterExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagFilterExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagFilterExprContext)
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




func (p *SQLParser) FuncParam() (localctx IFuncParamContext) {
	localctx = NewFuncParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, SQLParserRULE_funcParam)

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

	p.SetState(348)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 32, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(346)
			p.fieldExpr(0)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(347)
			p.tagFilterExpr(0)
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




func (p *SQLParser) ExprAtom() (localctx IExprAtomContext) {
	localctx = NewExprAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, SQLParserRULE_exprAtom)

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

	p.SetState(356)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(350)
			p.Ident()
		}
		p.SetState(352)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(351)
				p.IdentFilter()
			}


		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(354)
			p.DecNumber()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(355)
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

func (s *IdentFilterContext) TagFilterExpr() ITagFilterExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagFilterExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagFilterExprContext)
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




func (p *SQLParser) IdentFilter() (localctx IIdentFilterContext) {
	localctx = NewIdentFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, SQLParserRULE_identFilter)

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
		p.SetState(358)
		p.Match(SQLParserT_OPEN_SB)
	}
	{
		p.SetState(359)
		p.tagFilterExpr(0)
	}
	{
		p.SetState(360)
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




func (p *SQLParser) IntNumber() (localctx IIntNumberContext) {
	localctx = NewIntNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, SQLParserRULE_intNumber)
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
	p.SetState(363)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(362)
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
		p.SetState(365)
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




func (p *SQLParser) DecNumber() (localctx IDecNumberContext) {
	localctx = NewDecNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, SQLParserRULE_decNumber)
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
	p.SetState(368)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(367)
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
		p.SetState(370)
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




func (p *SQLParser) LimitClause() (localctx ILimitClauseContext) {
	localctx = NewLimitClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, SQLParserRULE_limitClause)

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
		p.SetState(372)
		p.Match(SQLParserT_LIMIT)
	}
	{
		p.SetState(373)
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




func (p *SQLParser) MetricName() (localctx IMetricNameContext) {
	localctx = NewMetricNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, SQLParserRULE_metricName)

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
		p.SetState(375)
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




func (p *SQLParser) TagKey() (localctx ITagKeyContext) {
	localctx = NewTagKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, SQLParserRULE_tagKey)

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
		p.SetState(377)
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




func (p *SQLParser) TagValue() (localctx ITagValueContext) {
	localctx = NewTagValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, SQLParserRULE_tagValue)

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
		p.SetState(379)
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




func (p *SQLParser) Ident() (localctx IIdentContext) {
	localctx = NewIdentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, SQLParserRULE_ident)

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
	p.SetState(383)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserL_ID:
		{
			p.SetState(381)
			p.Match(SQLParserL_ID)
		}


	case SQLParserT_CREATE, SQLParserT_INTERVAL, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NODE, SQLParserT_MEASUREMENTS, SQLParserT_MEASUREMENT, SQLParserT_FIELD, SQLParserT_TAG, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_PROFILE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
		{
			p.SetState(382)
			p.NonReservedWords()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(392)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(385)
				p.Match(SQLParserT_DOT)
			}
			p.SetState(388)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case SQLParserL_ID:
				{
					p.SetState(386)
					p.Match(SQLParserL_ID)
				}


			case SQLParserT_CREATE, SQLParserT_INTERVAL, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NODE, SQLParserT_MEASUREMENTS, SQLParserT_MEASUREMENT, SQLParserT_FIELD, SQLParserT_TAG, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_PROFILE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
				{
					p.SetState(387)
					p.NonReservedWords()
				}



			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}


		}
		p.SetState(394)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext())
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




func (p *SQLParser) NonReservedWords() (localctx INonReservedWordsContext) {
	localctx = NewNonReservedWordsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, SQLParserRULE_nonReservedWords)
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
		p.SetState(395)
		_la = p.GetTokenStream().LA(1)

		if !((((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_DATASBAE) | (1 << SQLParserT_DATASBAES) | (1 << SQLParserT_NODE) | (1 << SQLParserT_MEASUREMENTS) | (1 << SQLParserT_MEASUREMENT) | (1 << SQLParserT_FIELD) | (1 << SQLParserT_TAG) | (1 << SQLParserT_KEYS) | (1 << SQLParserT_KEY) | (1 << SQLParserT_WITH) | (1 << SQLParserT_VALUES) | (1 << SQLParserT_FROM) | (1 << SQLParserT_WHERE) | (1 << SQLParserT_LIMIT))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)) | (1 << (SQLParserT_ASC - 32)) | (1 << (SQLParserT_DESC - 32)) | (1 << (SQLParserT_LIKE - 32)) | (1 << (SQLParserT_NOT - 32)) | (1 << (SQLParserT_BETWEEN - 32)) | (1 << (SQLParserT_IS - 32)) | (1 << (SQLParserT_GROUP - 32)) | (1 << (SQLParserT_BY - 32)) | (1 << (SQLParserT_FOR - 32)) | (1 << (SQLParserT_STATS - 32)) | (1 << (SQLParserT_TIME - 32)) | (1 << (SQLParserT_PROFILE - 32)) | (1 << (SQLParserT_SECOND - 32)) | (1 << (SQLParserT_MINUTE - 32)) | (1 << (SQLParserT_HOUR - 32)) | (1 << (SQLParserT_DAY - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)))) != 0)) {
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
	case 10:
			var t *TagFilterExprContext = nil
			if localctx != nil { t = localctx.(*TagFilterExprContext) }
			return p.TagFilterExpr_Sempred(t, predIndex)

	case 24:
			var t *BoolExprContext = nil
			if localctx != nil { t = localctx.(*BoolExprContext) }
			return p.BoolExpr_Sempred(t, predIndex)

	case 29:
			var t *FieldExprContext = nil
			if localctx != nil { t = localctx.(*FieldExprContext) }
			return p.FieldExpr_Sempred(t, predIndex)


	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *SQLParser) TagFilterExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
			return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) BoolExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 1:
			return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) FieldExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
			return p.Precpred(p.GetParserRuleContext(), 8)

	case 3:
			return p.Precpred(p.GetParserRuleContext(), 7)

	case 4:
			return p.Precpred(p.GetParserRuleContext(), 6)

	case 5:
			return p.Precpred(p.GetParserRuleContext(), 5)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

