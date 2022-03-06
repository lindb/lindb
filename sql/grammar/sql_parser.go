// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated from /Users/jacklhuang/Documents/code/gopath/src/github.com/lindb/lindb/sql/grammar/SQL.g4 by ANTLR 4.9.2. DO NOT EDIT.

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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 106, 517, 
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
	9, 55, 4, 56, 9, 56, 4, 57, 9, 57, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 126, 10, 3, 3, 4, 3, 4, 3, 4, 3, 5, 
	3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 5, 6, 140, 10, 6, 3, 6, 
	5, 6, 143, 10, 6, 3, 7, 3, 7, 3, 7, 3, 7, 5, 7, 149, 10, 7, 3, 7, 3, 7, 
	3, 7, 3, 7, 5, 7, 155, 10, 7, 3, 7, 5, 7, 158, 10, 7, 3, 8, 3, 8, 3, 8, 
	3, 8, 5, 8, 164, 10, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 5, 9, 
	173, 10, 9, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 5, 10, 182, 
	10, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 5, 10, 190, 10, 10, 3, 
	10, 5, 10, 193, 10, 10, 3, 11, 3, 11, 3, 12, 3, 12, 3, 13, 3, 13, 3, 14, 
	5, 14, 202, 10, 14, 3, 14, 3, 14, 3, 14, 5, 14, 207, 10, 14, 3, 14, 3, 
	14, 5, 14, 211, 10, 14, 3, 14, 5, 14, 214, 10, 14, 3, 14, 5, 14, 217, 10, 
	14, 3, 14, 5, 14, 220, 10, 14, 3, 14, 5, 14, 223, 10, 14, 3, 15, 3, 15, 
	3, 15, 3, 16, 3, 16, 3, 16, 7, 16, 231, 10, 16, 12, 16, 14, 16, 234, 11, 
	16, 3, 17, 3, 17, 5, 17, 238, 10, 17, 3, 18, 3, 18, 3, 18, 3, 19, 3, 19, 
	3, 19, 3, 20, 3, 20, 3, 20, 3, 21, 3, 21, 3, 21, 3, 21, 3, 21, 3, 21, 3, 
	21, 3, 21, 5, 21, 257, 10, 21, 5, 21, 259, 10, 21, 3, 22, 3, 22, 3, 22, 
	3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 
	22, 5, 22, 275, 10, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 5, 22, 
	283, 10, 22, 3, 22, 3, 22, 3, 22, 3, 22, 5, 22, 289, 10, 22, 3, 22, 3, 
	22, 3, 22, 7, 22, 294, 10, 22, 12, 22, 14, 22, 297, 11, 22, 3, 23, 3, 23, 
	3, 23, 7, 23, 302, 10, 23, 12, 23, 14, 23, 305, 11, 23, 3, 24, 3, 24, 3, 
	24, 5, 24, 310, 10, 24, 3, 25, 3, 25, 3, 25, 3, 25, 5, 25, 316, 10, 25, 
	3, 26, 3, 26, 5, 26, 320, 10, 26, 3, 27, 3, 27, 3, 27, 5, 27, 325, 10, 
	27, 3, 27, 3, 27, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 
	5, 28, 337, 10, 28, 3, 28, 5, 28, 340, 10, 28, 3, 29, 3, 29, 3, 29, 7, 
	29, 345, 10, 29, 12, 29, 14, 29, 348, 11, 29, 3, 30, 3, 30, 3, 30, 3, 30, 
	3, 30, 3, 30, 5, 30, 356, 10, 30, 3, 31, 3, 31, 3, 32, 3, 32, 3, 32, 3, 
	32, 3, 33, 3, 33, 7, 33, 366, 10, 33, 12, 33, 14, 33, 369, 11, 33, 3, 34, 
	3, 34, 3, 34, 7, 34, 374, 10, 34, 12, 34, 14, 34, 377, 11, 34, 3, 35, 3, 
	35, 3, 35, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36, 5, 36, 388, 10, 36, 
	3, 36, 3, 36, 3, 36, 3, 36, 7, 36, 394, 10, 36, 12, 36, 14, 36, 397, 11, 
	36, 3, 37, 3, 37, 3, 38, 3, 38, 3, 39, 3, 39, 3, 39, 3, 39, 3, 40, 3, 40, 
	3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 5, 40, 415, 10, 40, 3, 41, 3, 
	41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 5, 41, 425, 10, 41, 3, 41, 
	3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 
	41, 7, 41, 439, 10, 41, 12, 41, 14, 41, 442, 11, 41, 3, 42, 3, 42, 3, 42, 
	3, 43, 3, 43, 3, 44, 3, 44, 3, 44, 5, 44, 452, 10, 44, 3, 44, 3, 44, 3, 
	45, 3, 45, 3, 46, 3, 46, 3, 46, 7, 46, 461, 10, 46, 12, 46, 14, 46, 464, 
	11, 46, 3, 47, 3, 47, 5, 47, 468, 10, 47, 3, 48, 3, 48, 5, 48, 472, 10, 
	48, 3, 48, 3, 48, 5, 48, 476, 10, 48, 3, 49, 3, 49, 3, 49, 3, 49, 3, 50, 
	5, 50, 483, 10, 50, 3, 50, 3, 50, 3, 51, 5, 51, 488, 10, 51, 3, 51, 3, 
	51, 3, 52, 3, 52, 3, 52, 3, 53, 3, 53, 3, 54, 3, 54, 3, 55, 3, 55, 3, 56, 
	3, 56, 5, 56, 503, 10, 56, 3, 56, 3, 56, 3, 56, 5, 56, 508, 10, 56, 7, 
	56, 510, 10, 56, 12, 56, 14, 56, 513, 11, 56, 3, 57, 3, 57, 3, 57, 2, 5, 
	42, 70, 80, 58, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 
	32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 
	68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 
	104, 106, 108, 110, 112, 2, 10, 3, 2, 44, 45, 4, 2, 47, 48, 104, 105, 3, 
	2, 50, 51, 4, 2, 52, 52, 89, 89, 3, 2, 73, 79, 3, 2, 66, 72, 3, 2, 98, 
	99, 4, 2, 3, 17, 19, 79, 2, 537, 2, 114, 3, 2, 2, 2, 4, 125, 3, 2, 2, 2, 
	6, 127, 3, 2, 2, 2, 8, 130, 3, 2, 2, 2, 10, 133, 3, 2, 2, 2, 12, 144, 3, 
	2, 2, 2, 14, 159, 3, 2, 2, 2, 16, 167, 3, 2, 2, 2, 18, 176, 3, 2, 2, 2, 
	20, 194, 3, 2, 2, 2, 22, 196, 3, 2, 2, 2, 24, 198, 3, 2, 2, 2, 26, 201, 
	3, 2, 2, 2, 28, 224, 3, 2, 2, 2, 30, 227, 3, 2, 2, 2, 32, 235, 3, 2, 2, 
	2, 34, 239, 3, 2, 2, 2, 36, 242, 3, 2, 2, 2, 38, 245, 3, 2, 2, 2, 40, 258, 
	3, 2, 2, 2, 42, 288, 3, 2, 2, 2, 44, 298, 3, 2, 2, 2, 46, 306, 3, 2, 2, 
	2, 48, 311, 3, 2, 2, 2, 50, 317, 3, 2, 2, 2, 52, 321, 3, 2, 2, 2, 54, 328, 
	3, 2, 2, 2, 56, 341, 3, 2, 2, 2, 58, 355, 3, 2, 2, 2, 60, 357, 3, 2, 2, 
	2, 62, 359, 3, 2, 2, 2, 64, 363, 3, 2, 2, 2, 66, 370, 3, 2, 2, 2, 68, 378, 
	3, 2, 2, 2, 70, 387, 3, 2, 2, 2, 72, 398, 3, 2, 2, 2, 74, 400, 3, 2, 2, 
	2, 76, 402, 3, 2, 2, 2, 78, 414, 3, 2, 2, 2, 80, 424, 3, 2, 2, 2, 82, 443, 
	3, 2, 2, 2, 84, 446, 3, 2, 2, 2, 86, 448, 3, 2, 2, 2, 88, 455, 3, 2, 2, 
	2, 90, 457, 3, 2, 2, 2, 92, 467, 3, 2, 2, 2, 94, 475, 3, 2, 2, 2, 96, 477, 
	3, 2, 2, 2, 98, 482, 3, 2, 2, 2, 100, 487, 3, 2, 2, 2, 102, 491, 3, 2, 
	2, 2, 104, 494, 3, 2, 2, 2, 106, 496, 3, 2, 2, 2, 108, 498, 3, 2, 2, 2, 
	110, 502, 3, 2, 2, 2, 112, 514, 3, 2, 2, 2, 114, 115, 5, 4, 3, 2, 115, 
	116, 7, 2, 2, 3, 116, 3, 3, 2, 2, 2, 117, 126, 5, 6, 4, 2, 118, 126, 5, 
	8, 5, 2, 119, 126, 5, 10, 6, 2, 120, 126, 5, 12, 7, 2, 121, 126, 5, 14, 
	8, 2, 122, 126, 5, 16, 9, 2, 123, 126, 5, 18, 10, 2, 124, 126, 5, 26, 14, 
	2, 125, 117, 3, 2, 2, 2, 125, 118, 3, 2, 2, 2, 125, 119, 3, 2, 2, 2, 125, 
	120, 3, 2, 2, 2, 125, 121, 3, 2, 2, 2, 125, 122, 3, 2, 2, 2, 125, 123, 
	3, 2, 2, 2, 125, 124, 3, 2, 2, 2, 126, 5, 3, 2, 2, 2, 127, 128, 7, 17, 
	2, 2, 128, 129, 7, 18, 2, 2, 129, 7, 3, 2, 2, 2, 130, 131, 7, 17, 2, 2, 
	131, 132, 7, 20, 2, 2, 132, 9, 3, 2, 2, 2, 133, 134, 7, 17, 2, 2, 134, 
	139, 7, 22, 2, 2, 135, 136, 7, 36, 2, 2, 136, 137, 7, 21, 2, 2, 137, 138, 
	7, 82, 2, 2, 138, 140, 5, 20, 11, 2, 139, 135, 3, 2, 2, 2, 139, 140, 3, 
	2, 2, 2, 140, 142, 3, 2, 2, 2, 141, 143, 5, 102, 52, 2, 142, 141, 3, 2, 
	2, 2, 142, 143, 3, 2, 2, 2, 143, 11, 3, 2, 2, 2, 144, 145, 7, 17, 2, 2, 
	145, 148, 7, 24, 2, 2, 146, 147, 7, 16, 2, 2, 147, 149, 5, 24, 13, 2, 148, 
	146, 3, 2, 2, 2, 148, 149, 3, 2, 2, 2, 149, 154, 3, 2, 2, 2, 150, 151, 
	7, 36, 2, 2, 151, 152, 7, 25, 2, 2, 152, 153, 7, 82, 2, 2, 153, 155, 5, 
	20, 11, 2, 154, 150, 3, 2, 2, 2, 154, 155, 3, 2, 2, 2, 155, 157, 3, 2, 
	2, 2, 156, 158, 5, 102, 52, 2, 157, 156, 3, 2, 2, 2, 157, 158, 3, 2, 2, 
	2, 158, 13, 3, 2, 2, 2, 159, 160, 7, 17, 2, 2, 160, 163, 7, 27, 2, 2, 161, 
	162, 7, 16, 2, 2, 162, 164, 5, 24, 13, 2, 163, 161, 3, 2, 2, 2, 163, 164, 
	3, 2, 2, 2, 164, 165, 3, 2, 2, 2, 165, 166, 5, 36, 19, 2, 166, 15, 3, 2, 
	2, 2, 167, 168, 7, 17, 2, 2, 168, 169, 7, 28, 2, 2, 169, 172, 7, 30, 2, 
	2, 170, 171, 7, 16, 2, 2, 171, 173, 5, 24, 13, 2, 172, 170, 3, 2, 2, 2, 
	172, 173, 3, 2, 2, 2, 173, 174, 3, 2, 2, 2, 174, 175, 5, 36, 19, 2, 175, 
	17, 3, 2, 2, 2, 176, 177, 7, 17, 2, 2, 177, 178, 7, 28, 2, 2, 178, 181, 
	7, 33, 2, 2, 179, 180, 7, 16, 2, 2, 180, 182, 5, 24, 13, 2, 181, 179, 3, 
	2, 2, 2, 181, 182, 3, 2, 2, 2, 182, 183, 3, 2, 2, 2, 183, 184, 5, 36, 19, 
	2, 184, 185, 7, 32, 2, 2, 185, 186, 7, 31, 2, 2, 186, 187, 7, 82, 2, 2, 
	187, 189, 5, 22, 12, 2, 188, 190, 5, 38, 20, 2, 189, 188, 3, 2, 2, 2, 189, 
	190, 3, 2, 2, 2, 190, 192, 3, 2, 2, 2, 191, 193, 5, 102, 52, 2, 192, 191, 
	3, 2, 2, 2, 192, 193, 3, 2, 2, 2, 193, 19, 3, 2, 2, 2, 194, 195, 5, 110, 
	56, 2, 195, 21, 3, 2, 2, 2, 196, 197, 5, 110, 56, 2, 197, 23, 3, 2, 2, 
	2, 198, 199, 5, 110, 56, 2, 199, 25, 3, 2, 2, 2, 200, 202, 7, 40, 2, 2, 
	201, 200, 3, 2, 2, 2, 201, 202, 3, 2, 2, 2, 202, 203, 3, 2, 2, 2, 203, 
	206, 5, 28, 15, 2, 204, 205, 7, 16, 2, 2, 205, 207, 5, 24, 13, 2, 206, 
	204, 3, 2, 2, 2, 206, 207, 3, 2, 2, 2, 207, 208, 3, 2, 2, 2, 208, 210, 
	5, 36, 19, 2, 209, 211, 5, 38, 20, 2, 210, 209, 3, 2, 2, 2, 210, 211, 3, 
	2, 2, 2, 211, 213, 3, 2, 2, 2, 212, 214, 5, 54, 28, 2, 213, 212, 3, 2, 
	2, 2, 213, 214, 3, 2, 2, 2, 214, 216, 3, 2, 2, 2, 215, 217, 5, 62, 32, 
	2, 216, 215, 3, 2, 2, 2, 216, 217, 3, 2, 2, 2, 217, 219, 3, 2, 2, 2, 218, 
	220, 5, 102, 52, 2, 219, 218, 3, 2, 2, 2, 219, 220, 3, 2, 2, 2, 220, 222, 
	3, 2, 2, 2, 221, 223, 7, 41, 2, 2, 222, 221, 3, 2, 2, 2, 222, 223, 3, 2, 
	2, 2, 223, 27, 3, 2, 2, 2, 224, 225, 7, 42, 2, 2, 225, 226, 5, 30, 16, 
	2, 226, 29, 3, 2, 2, 2, 227, 232, 5, 32, 17, 2, 228, 229, 7, 91, 2, 2, 
	229, 231, 5, 32, 17, 2, 230, 228, 3, 2, 2, 2, 231, 234, 3, 2, 2, 2, 232, 
	230, 3, 2, 2, 2, 232, 233, 3, 2, 2, 2, 233, 31, 3, 2, 2, 2, 234, 232, 3, 
	2, 2, 2, 235, 237, 5, 80, 41, 2, 236, 238, 5, 34, 18, 2, 237, 236, 3, 2, 
	2, 2, 237, 238, 3, 2, 2, 2, 238, 33, 3, 2, 2, 2, 239, 240, 7, 43, 2, 2, 
	240, 241, 5, 110, 56, 2, 241, 35, 3, 2, 2, 2, 242, 243, 7, 35, 2, 2, 243, 
	244, 5, 104, 53, 2, 244, 37, 3, 2, 2, 2, 245, 246, 7, 36, 2, 2, 246, 247, 
	5, 40, 21, 2, 247, 39, 3, 2, 2, 2, 248, 259, 5, 42, 22, 2, 249, 250, 5, 
	42, 22, 2, 250, 251, 7, 44, 2, 2, 251, 252, 5, 46, 24, 2, 252, 259, 3, 
	2, 2, 2, 253, 256, 5, 46, 24, 2, 254, 255, 7, 44, 2, 2, 255, 257, 5, 42, 
	22, 2, 256, 254, 3, 2, 2, 2, 256, 257, 3, 2, 2, 2, 257, 259, 3, 2, 2, 2, 
	258, 248, 3, 2, 2, 2, 258, 249, 3, 2, 2, 2, 258, 253, 3, 2, 2, 2, 259, 
	41, 3, 2, 2, 2, 260, 261, 8, 22, 1, 2, 261, 262, 7, 96, 2, 2, 262, 263, 
	5, 42, 22, 2, 263, 264, 7, 97, 2, 2, 264, 289, 3, 2, 2, 2, 265, 274, 5, 
	106, 54, 2, 266, 275, 7, 82, 2, 2, 267, 275, 7, 52, 2, 2, 268, 269, 7, 
	53, 2, 2, 269, 275, 7, 52, 2, 2, 270, 275, 7, 89, 2, 2, 271, 275, 7, 90, 
	2, 2, 272, 275, 7, 83, 2, 2, 273, 275, 7, 84, 2, 2, 274, 266, 3, 2, 2, 
	2, 274, 267, 3, 2, 2, 2, 274, 268, 3, 2, 2, 2, 274, 270, 3, 2, 2, 2, 274, 
	271, 3, 2, 2, 2, 274, 272, 3, 2, 2, 2, 274, 273, 3, 2, 2, 2, 275, 276, 
	3, 2, 2, 2, 276, 277, 5, 108, 55, 2, 277, 289, 3, 2, 2, 2, 278, 282, 5, 
	106, 54, 2, 279, 283, 7, 63, 2, 2, 280, 281, 7, 53, 2, 2, 281, 283, 7, 
	63, 2, 2, 282, 279, 3, 2, 2, 2, 282, 280, 3, 2, 2, 2, 283, 284, 3, 2, 2, 
	2, 284, 285, 7, 96, 2, 2, 285, 286, 5, 44, 23, 2, 286, 287, 7, 97, 2, 2, 
	287, 289, 3, 2, 2, 2, 288, 260, 3, 2, 2, 2, 288, 265, 3, 2, 2, 2, 288, 
	278, 3, 2, 2, 2, 289, 295, 3, 2, 2, 2, 290, 291, 12, 3, 2, 2, 291, 292, 
	9, 2, 2, 2, 292, 294, 5, 42, 22, 4, 293, 290, 3, 2, 2, 2, 294, 297, 3, 
	2, 2, 2, 295, 293, 3, 2, 2, 2, 295, 296, 3, 2, 2, 2, 296, 43, 3, 2, 2, 
	2, 297, 295, 3, 2, 2, 2, 298, 303, 5, 108, 55, 2, 299, 300, 7, 91, 2, 2, 
	300, 302, 5, 108, 55, 2, 301, 299, 3, 2, 2, 2, 302, 305, 3, 2, 2, 2, 303, 
	301, 3, 2, 2, 2, 303, 304, 3, 2, 2, 2, 304, 45, 3, 2, 2, 2, 305, 303, 3, 
	2, 2, 2, 306, 309, 5, 48, 25, 2, 307, 308, 7, 44, 2, 2, 308, 310, 5, 48, 
	25, 2, 309, 307, 3, 2, 2, 2, 309, 310, 3, 2, 2, 2, 310, 47, 3, 2, 2, 2, 
	311, 312, 7, 61, 2, 2, 312, 315, 5, 78, 40, 2, 313, 316, 5, 50, 26, 2, 
	314, 316, 5, 110, 56, 2, 315, 313, 3, 2, 2, 2, 315, 314, 3, 2, 2, 2, 316, 
	49, 3, 2, 2, 2, 317, 319, 5, 52, 27, 2, 318, 320, 5, 82, 42, 2, 319, 318, 
	3, 2, 2, 2, 319, 320, 3, 2, 2, 2, 320, 51, 3, 2, 2, 2, 321, 322, 7, 62, 
	2, 2, 322, 324, 7, 96, 2, 2, 323, 325, 5, 90, 46, 2, 324, 323, 3, 2, 2, 
	2, 324, 325, 3, 2, 2, 2, 325, 326, 3, 2, 2, 2, 326, 327, 7, 97, 2, 2, 327, 
	53, 3, 2, 2, 2, 328, 329, 7, 56, 2, 2, 329, 330, 7, 58, 2, 2, 330, 336, 
	5, 56, 29, 2, 331, 332, 7, 46, 2, 2, 332, 333, 7, 96, 2, 2, 333, 334, 5, 
	60, 31, 2, 334, 335, 7, 97, 2, 2, 335, 337, 3, 2, 2, 2, 336, 331, 3, 2, 
	2, 2, 336, 337, 3, 2, 2, 2, 337, 339, 3, 2, 2, 2, 338, 340, 5, 68, 35, 
	2, 339, 338, 3, 2, 2, 2, 339, 340, 3, 2, 2, 2, 340, 55, 3, 2, 2, 2, 341, 
	346, 5, 58, 30, 2, 342, 343, 7, 91, 2, 2, 343, 345, 5, 58, 30, 2, 344, 
	342, 3, 2, 2, 2, 345, 348, 3, 2, 2, 2, 346, 344, 3, 2, 2, 2, 346, 347, 
	3, 2, 2, 2, 347, 57, 3, 2, 2, 2, 348, 346, 3, 2, 2, 2, 349, 356, 5, 110, 
	56, 2, 350, 351, 7, 61, 2, 2, 351, 352, 7, 96, 2, 2, 352, 353, 5, 82, 42, 
	2, 353, 354, 7, 97, 2, 2, 354, 356, 3, 2, 2, 2, 355, 349, 3, 2, 2, 2, 355, 
	350, 3, 2, 2, 2, 356, 59, 3, 2, 2, 2, 357, 358, 9, 3, 2, 2, 358, 61, 3, 
	2, 2, 2, 359, 360, 7, 49, 2, 2, 360, 361, 7, 58, 2, 2, 361, 362, 5, 66, 
	34, 2, 362, 63, 3, 2, 2, 2, 363, 367, 5, 80, 41, 2, 364, 366, 9, 4, 2, 
	2, 365, 364, 3, 2, 2, 2, 366, 369, 3, 2, 2, 2, 367, 365, 3, 2, 2, 2, 367, 
	368, 3, 2, 2, 2, 368, 65, 3, 2, 2, 2, 369, 367, 3, 2, 2, 2, 370, 375, 5, 
	64, 33, 2, 371, 372, 7, 91, 2, 2, 372, 374, 5, 64, 33, 2, 373, 371, 3, 
	2, 2, 2, 374, 377, 3, 2, 2, 2, 375, 373, 3, 2, 2, 2, 375, 376, 3, 2, 2, 
	2, 376, 67, 3, 2, 2, 2, 377, 375, 3, 2, 2, 2, 378, 379, 7, 57, 2, 2, 379, 
	380, 5, 70, 36, 2, 380, 69, 3, 2, 2, 2, 381, 382, 8, 36, 1, 2, 382, 383, 
	7, 96, 2, 2, 383, 384, 5, 70, 36, 2, 384, 385, 7, 97, 2, 2, 385, 388, 3, 
	2, 2, 2, 386, 388, 5, 74, 38, 2, 387, 381, 3, 2, 2, 2, 387, 386, 3, 2, 
	2, 2, 388, 395, 3, 2, 2, 2, 389, 390, 12, 4, 2, 2, 390, 391, 5, 72, 37, 
	2, 391, 392, 5, 70, 36, 5, 392, 394, 3, 2, 2, 2, 393, 389, 3, 2, 2, 2, 
	394, 397, 3, 2, 2, 2, 395, 393, 3, 2, 2, 2, 395, 396, 3, 2, 2, 2, 396, 
	71, 3, 2, 2, 2, 397, 395, 3, 2, 2, 2, 398, 399, 9, 2, 2, 2, 399, 73, 3, 
	2, 2, 2, 400, 401, 5, 76, 39, 2, 401, 75, 3, 2, 2, 2, 402, 403, 5, 80, 
	41, 2, 403, 404, 5, 78, 40, 2, 404, 405, 5, 80, 41, 2, 405, 77, 3, 2, 2, 
	2, 406, 415, 7, 82, 2, 2, 407, 415, 7, 83, 2, 2, 408, 415, 7, 84, 2, 2, 
	409, 415, 7, 87, 2, 2, 410, 415, 7, 88, 2, 2, 411, 415, 7, 85, 2, 2, 412, 
	415, 7, 86, 2, 2, 413, 415, 9, 5, 2, 2, 414, 406, 3, 2, 2, 2, 414, 407, 
	3, 2, 2, 2, 414, 408, 3, 2, 2, 2, 414, 409, 3, 2, 2, 2, 414, 410, 3, 2, 
	2, 2, 414, 411, 3, 2, 2, 2, 414, 412, 3, 2, 2, 2, 414, 413, 3, 2, 2, 2, 
	415, 79, 3, 2, 2, 2, 416, 417, 8, 41, 1, 2, 417, 418, 7, 96, 2, 2, 418, 
	419, 5, 80, 41, 2, 419, 420, 7, 97, 2, 2, 420, 425, 3, 2, 2, 2, 421, 425, 
	5, 86, 44, 2, 422, 425, 5, 94, 48, 2, 423, 425, 5, 82, 42, 2, 424, 416, 
	3, 2, 2, 2, 424, 421, 3, 2, 2, 2, 424, 422, 3, 2, 2, 2, 424, 423, 3, 2, 
	2, 2, 425, 440, 3, 2, 2, 2, 426, 427, 12, 10, 2, 2, 427, 428, 7, 101, 2, 
	2, 428, 439, 5, 80, 41, 11, 429, 430, 12, 9, 2, 2, 430, 431, 7, 100, 2, 
	2, 431, 439, 5, 80, 41, 10, 432, 433, 12, 8, 2, 2, 433, 434, 7, 98, 2, 
	2, 434, 439, 5, 80, 41, 9, 435, 436, 12, 7, 2, 2, 436, 437, 7, 99, 2, 2, 
	437, 439, 5, 80, 41, 8, 438, 426, 3, 2, 2, 2, 438, 429, 3, 2, 2, 2, 438, 
	432, 3, 2, 2, 2, 438, 435, 3, 2, 2, 2, 439, 442, 3, 2, 2, 2, 440, 438, 
	3, 2, 2, 2, 440, 441, 3, 2, 2, 2, 441, 81, 3, 2, 2, 2, 442, 440, 3, 2, 
	2, 2, 443, 444, 5, 98, 50, 2, 444, 445, 5, 84, 43, 2, 445, 83, 3, 2, 2, 
	2, 446, 447, 9, 6, 2, 2, 447, 85, 3, 2, 2, 2, 448, 449, 5, 88, 45, 2, 449, 
	451, 7, 96, 2, 2, 450, 452, 5, 90, 46, 2, 451, 450, 3, 2, 2, 2, 451, 452, 
	3, 2, 2, 2, 452, 453, 3, 2, 2, 2, 453, 454, 7, 97, 2, 2, 454, 87, 3, 2, 
	2, 2, 455, 456, 9, 7, 2, 2, 456, 89, 3, 2, 2, 2, 457, 462, 5, 92, 47, 2, 
	458, 459, 7, 91, 2, 2, 459, 461, 5, 92, 47, 2, 460, 458, 3, 2, 2, 2, 461, 
	464, 3, 2, 2, 2, 462, 460, 3, 2, 2, 2, 462, 463, 3, 2, 2, 2, 463, 91, 3, 
	2, 2, 2, 464, 462, 3, 2, 2, 2, 465, 468, 5, 80, 41, 2, 466, 468, 5, 42, 
	22, 2, 467, 465, 3, 2, 2, 2, 467, 466, 3, 2, 2, 2, 468, 93, 3, 2, 2, 2, 
	469, 471, 5, 110, 56, 2, 470, 472, 5, 96, 49, 2, 471, 470, 3, 2, 2, 2, 
	471, 472, 3, 2, 2, 2, 472, 476, 3, 2, 2, 2, 473, 476, 5, 100, 51, 2, 474, 
	476, 5, 98, 50, 2, 475, 469, 3, 2, 2, 2, 475, 473, 3, 2, 2, 2, 475, 474, 
	3, 2, 2, 2, 476, 95, 3, 2, 2, 2, 477, 478, 7, 94, 2, 2, 478, 479, 5, 42, 
	22, 2, 479, 480, 7, 95, 2, 2, 480, 97, 3, 2, 2, 2, 481, 483, 9, 8, 2, 2, 
	482, 481, 3, 2, 2, 2, 482, 483, 3, 2, 2, 2, 483, 484, 3, 2, 2, 2, 484, 
	485, 7, 104, 2, 2, 485, 99, 3, 2, 2, 2, 486, 488, 9, 8, 2, 2, 487, 486, 
	3, 2, 2, 2, 487, 488, 3, 2, 2, 2, 488, 489, 3, 2, 2, 2, 489, 490, 7, 105, 
	2, 2, 490, 101, 3, 2, 2, 2, 491, 492, 7, 37, 2, 2, 492, 493, 7, 104, 2, 
	2, 493, 103, 3, 2, 2, 2, 494, 495, 5, 110, 56, 2, 495, 105, 3, 2, 2, 2, 
	496, 497, 5, 110, 56, 2, 497, 107, 3, 2, 2, 2, 498, 499, 5, 110, 56, 2, 
	499, 109, 3, 2, 2, 2, 500, 503, 7, 103, 2, 2, 501, 503, 5, 112, 57, 2, 
	502, 500, 3, 2, 2, 2, 502, 501, 3, 2, 2, 2, 503, 511, 3, 2, 2, 2, 504, 
	507, 7, 80, 2, 2, 505, 508, 7, 103, 2, 2, 506, 508, 5, 112, 57, 2, 507, 
	505, 3, 2, 2, 2, 507, 506, 3, 2, 2, 2, 508, 510, 3, 2, 2, 2, 509, 504, 
	3, 2, 2, 2, 510, 513, 3, 2, 2, 2, 511, 509, 3, 2, 2, 2, 511, 512, 3, 2, 
	2, 2, 512, 111, 3, 2, 2, 2, 513, 511, 3, 2, 2, 2, 514, 515, 9, 9, 2, 2, 
	515, 113, 3, 2, 2, 2, 55, 125, 139, 142, 148, 154, 157, 163, 172, 181, 
	189, 192, 201, 206, 210, 213, 216, 219, 222, 232, 237, 256, 258, 274, 282, 
	288, 295, 303, 309, 315, 319, 324, 336, 339, 346, 355, 367, 375, 387, 395, 
	414, 424, 438, 440, 451, 462, 467, 471, 475, 482, 487, 502, 507, 511,
}
var literalNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"'m'", "", "", "", "'M'", "", "'.'", "':'", "'='", "'<>'", "'!='", "'>'", 
	"'>='", "'<'", "'<='", "'=~'", "'!~'", "','", "'{'", "'}'", "'['", "']'", 
	"'('", "')'", "'+'", "'-'", "'/'", "'*'", "'%'",
}
var symbolicNames = []string{
	"", "T_CREATE", "T_UPDATE", "T_SET", "T_DROP", "T_INTERVAL", "T_INTERVAL_NAME", 
	"T_SHARD", "T_REPLICATION", "T_TTL", "T_META_TTL", "T_PAST_TTL", "T_FUTURE_TTL", 
	"T_KILL", "T_ON", "T_SHOW", "T_MASTER", "T_DATASBAE", "T_DATASBAES", "T_NAMESPACE", 
	"T_NAMESPACES", "T_NODE", "T_METRICS", "T_METRIC", "T_FIELD", "T_FIELDS", 
	"T_TAG", "T_INFO", "T_KEYS", "T_KEY", "T_WITH", "T_VALUES", "T_VALUE", 
	"T_FROM", "T_WHERE", "T_LIMIT", "T_QUERIES", "T_QUERY", "T_EXPLAIN", "T_WITH_VALUE", 
	"T_SELECT", "T_AS", "T_AND", "T_OR", "T_FILL", "T_NULL", "T_PREVIOUS", 
	"T_ORDER", "T_ASC", "T_DESC", "T_LIKE", "T_NOT", "T_BETWEEN", "T_IS", "T_GROUP", 
	"T_HAVING", "T_BY", "T_FOR", "T_STATS", "T_TIME", "T_NOW", "T_IN", "T_LOG", 
	"T_PROFILE", "T_SUM", "T_MIN", "T_MAX", "T_COUNT", "T_AVG", "T_STDDEV", 
	"T_QUANTILE", "T_SECOND", "T_MINUTE", "T_HOUR", "T_DAY", "T_WEEK", "T_MONTH", 
	"T_YEAR", "T_DOT", "T_COLON", "T_EQUAL", "T_NOTEQUAL", "T_NOTEQUAL2", "T_GREATER", 
	"T_GREATEREQUAL", "T_LESS", "T_LESSEQUAL", "T_REGEXP", "T_NEQREGEXP", "T_COMMA", 
	"T_OPEN_B", "T_CLOSE_B", "T_OPEN_SB", "T_CLOSE_SB", "T_OPEN_P", "T_CLOSE_P", 
	"T_ADD", "T_SUB", "T_DIV", "T_MUL", "T_MOD", "L_ID", "L_INT", "L_DEC", 
	"WS",
}

var ruleNames = []string{
	"statement", "statementList", "showMasterStmt", "showDatabaseStmt", "showNameSpacesStmt", 
	"showMetricsStmt", "showFieldsStmt", "showTagKeysStmt", "showTagValuesStmt", 
	"prefix", "withTagKey", "namespace", "queryStmt", "selectExpr", "fields", 
	"field", "alias", "fromClause", "whereClause", "conditionExpr", "tagFilterExpr", 
	"tagValueList", "timeRangeExpr", "timeExpr", "nowExpr", "nowFunc", "groupByClause", 
	"groupByKeys", "groupByKey", "fillOption", "orderByClause", "sortField", 
	"sortFields", "havingClause", "boolExpr", "boolExprLogicalOp", "boolExprAtom", 
	"binaryExpr", "binaryOperator", "fieldExpr", "durationLit", "intervalItem", 
	"exprFunc", "funcName", "exprFuncParams", "funcParam", "exprAtom", "identFilter", 
	"intNumber", "decNumber", "limitClause", "metricName", "tagKey", "tagValue", 
	"ident", "nonReservedWords",
}
type SQLParser struct {
	*antlr.BaseParser
}

// NewSQLParser produces a new parser instance for the optional input antlr.TokenStream.
//
// The *SQLParser instance produced may be reused by calling the SetInputStream method.
// The initial parser configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewSQLParser(input antlr.TokenStream) *SQLParser {
	this := new(SQLParser)
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
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
	SQLParserT_MASTER = 16
	SQLParserT_DATASBAE = 17
	SQLParserT_DATASBAES = 18
	SQLParserT_NAMESPACE = 19
	SQLParserT_NAMESPACES = 20
	SQLParserT_NODE = 21
	SQLParserT_METRICS = 22
	SQLParserT_METRIC = 23
	SQLParserT_FIELD = 24
	SQLParserT_FIELDS = 25
	SQLParserT_TAG = 26
	SQLParserT_INFO = 27
	SQLParserT_KEYS = 28
	SQLParserT_KEY = 29
	SQLParserT_WITH = 30
	SQLParserT_VALUES = 31
	SQLParserT_VALUE = 32
	SQLParserT_FROM = 33
	SQLParserT_WHERE = 34
	SQLParserT_LIMIT = 35
	SQLParserT_QUERIES = 36
	SQLParserT_QUERY = 37
	SQLParserT_EXPLAIN = 38
	SQLParserT_WITH_VALUE = 39
	SQLParserT_SELECT = 40
	SQLParserT_AS = 41
	SQLParserT_AND = 42
	SQLParserT_OR = 43
	SQLParserT_FILL = 44
	SQLParserT_NULL = 45
	SQLParserT_PREVIOUS = 46
	SQLParserT_ORDER = 47
	SQLParserT_ASC = 48
	SQLParserT_DESC = 49
	SQLParserT_LIKE = 50
	SQLParserT_NOT = 51
	SQLParserT_BETWEEN = 52
	SQLParserT_IS = 53
	SQLParserT_GROUP = 54
	SQLParserT_HAVING = 55
	SQLParserT_BY = 56
	SQLParserT_FOR = 57
	SQLParserT_STATS = 58
	SQLParserT_TIME = 59
	SQLParserT_NOW = 60
	SQLParserT_IN = 61
	SQLParserT_LOG = 62
	SQLParserT_PROFILE = 63
	SQLParserT_SUM = 64
	SQLParserT_MIN = 65
	SQLParserT_MAX = 66
	SQLParserT_COUNT = 67
	SQLParserT_AVG = 68
	SQLParserT_STDDEV = 69
	SQLParserT_QUANTILE = 70
	SQLParserT_SECOND = 71
	SQLParserT_MINUTE = 72
	SQLParserT_HOUR = 73
	SQLParserT_DAY = 74
	SQLParserT_WEEK = 75
	SQLParserT_MONTH = 76
	SQLParserT_YEAR = 77
	SQLParserT_DOT = 78
	SQLParserT_COLON = 79
	SQLParserT_EQUAL = 80
	SQLParserT_NOTEQUAL = 81
	SQLParserT_NOTEQUAL2 = 82
	SQLParserT_GREATER = 83
	SQLParserT_GREATEREQUAL = 84
	SQLParserT_LESS = 85
	SQLParserT_LESSEQUAL = 86
	SQLParserT_REGEXP = 87
	SQLParserT_NEQREGEXP = 88
	SQLParserT_COMMA = 89
	SQLParserT_OPEN_B = 90
	SQLParserT_CLOSE_B = 91
	SQLParserT_OPEN_SB = 92
	SQLParserT_CLOSE_SB = 93
	SQLParserT_OPEN_P = 94
	SQLParserT_CLOSE_P = 95
	SQLParserT_ADD = 96
	SQLParserT_SUB = 97
	SQLParserT_DIV = 98
	SQLParserT_MUL = 99
	SQLParserT_MOD = 100
	SQLParserL_ID = 101
	SQLParserL_INT = 102
	SQLParserL_DEC = 103
	SQLParserWS = 104
)

// SQLParser rules.
const (
	SQLParserRULE_statement = 0
	SQLParserRULE_statementList = 1
	SQLParserRULE_showMasterStmt = 2
	SQLParserRULE_showDatabaseStmt = 3
	SQLParserRULE_showNameSpacesStmt = 4
	SQLParserRULE_showMetricsStmt = 5
	SQLParserRULE_showFieldsStmt = 6
	SQLParserRULE_showTagKeysStmt = 7
	SQLParserRULE_showTagValuesStmt = 8
	SQLParserRULE_prefix = 9
	SQLParserRULE_withTagKey = 10
	SQLParserRULE_namespace = 11
	SQLParserRULE_queryStmt = 12
	SQLParserRULE_selectExpr = 13
	SQLParserRULE_fields = 14
	SQLParserRULE_field = 15
	SQLParserRULE_alias = 16
	SQLParserRULE_fromClause = 17
	SQLParserRULE_whereClause = 18
	SQLParserRULE_conditionExpr = 19
	SQLParserRULE_tagFilterExpr = 20
	SQLParserRULE_tagValueList = 21
	SQLParserRULE_timeRangeExpr = 22
	SQLParserRULE_timeExpr = 23
	SQLParserRULE_nowExpr = 24
	SQLParserRULE_nowFunc = 25
	SQLParserRULE_groupByClause = 26
	SQLParserRULE_groupByKeys = 27
	SQLParserRULE_groupByKey = 28
	SQLParserRULE_fillOption = 29
	SQLParserRULE_orderByClause = 30
	SQLParserRULE_sortField = 31
	SQLParserRULE_sortFields = 32
	SQLParserRULE_havingClause = 33
	SQLParserRULE_boolExpr = 34
	SQLParserRULE_boolExprLogicalOp = 35
	SQLParserRULE_boolExprAtom = 36
	SQLParserRULE_binaryExpr = 37
	SQLParserRULE_binaryOperator = 38
	SQLParserRULE_fieldExpr = 39
	SQLParserRULE_durationLit = 40
	SQLParserRULE_intervalItem = 41
	SQLParserRULE_exprFunc = 42
	SQLParserRULE_funcName = 43
	SQLParserRULE_exprFuncParams = 44
	SQLParserRULE_funcParam = 45
	SQLParserRULE_exprAtom = 46
	SQLParserRULE_identFilter = 47
	SQLParserRULE_intNumber = 48
	SQLParserRULE_decNumber = 49
	SQLParserRULE_limitClause = 50
	SQLParserRULE_metricName = 51
	SQLParserRULE_tagKey = 52
	SQLParserRULE_tagValue = 53
	SQLParserRULE_ident = 54
	SQLParserRULE_nonReservedWords = 55
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
		p.SetState(112)
		p.StatementList()
	}
	{
		p.SetState(113)
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

func (s *StatementListContext) ShowMasterStmt() IShowMasterStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowMasterStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowMasterStmtContext)
}

func (s *StatementListContext) ShowDatabaseStmt() IShowDatabaseStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowDatabaseStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowDatabaseStmtContext)
}

func (s *StatementListContext) ShowNameSpacesStmt() IShowNameSpacesStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowNameSpacesStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowNameSpacesStmtContext)
}

func (s *StatementListContext) ShowMetricsStmt() IShowMetricsStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowMetricsStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowMetricsStmtContext)
}

func (s *StatementListContext) ShowFieldsStmt() IShowFieldsStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowFieldsStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowFieldsStmtContext)
}

func (s *StatementListContext) ShowTagKeysStmt() IShowTagKeysStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowTagKeysStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowTagKeysStmtContext)
}

func (s *StatementListContext) ShowTagValuesStmt() IShowTagValuesStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowTagValuesStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowTagValuesStmtContext)
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

	p.SetState(123)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(115)
			p.ShowMasterStmt()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(116)
			p.ShowDatabaseStmt()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(117)
			p.ShowNameSpacesStmt()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(118)
			p.ShowMetricsStmt()
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(119)
			p.ShowFieldsStmt()
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(120)
			p.ShowTagKeysStmt()
		}


	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(121)
			p.ShowTagValuesStmt()
		}


	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(122)
			p.QueryStmt()
		}

	}


	return localctx
}


// IShowMasterStmtContext is an interface to support dynamic dispatch.
type IShowMasterStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowMasterStmtContext differentiates from other interfaces.
	IsShowMasterStmtContext()
}

type ShowMasterStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowMasterStmtContext() *ShowMasterStmtContext {
	var p = new(ShowMasterStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showMasterStmt
	return p
}

func (*ShowMasterStmtContext) IsShowMasterStmtContext() {}

func NewShowMasterStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowMasterStmtContext {
	var p = new(ShowMasterStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showMasterStmt

	return p
}

func (s *ShowMasterStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowMasterStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowMasterStmtContext) T_MASTER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MASTER, 0)
}

func (s *ShowMasterStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMasterStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowMasterStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowMasterStmt(s)
	}
}

func (s *ShowMasterStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowMasterStmt(s)
	}
}




func (p *SQLParser) ShowMasterStmt() (localctx IShowMasterStmtContext) {
	localctx = NewShowMasterStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SQLParserRULE_showMasterStmt)

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
		p.SetState(125)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(126)
		p.Match(SQLParserT_MASTER)
	}



	return localctx
}


// IShowDatabaseStmtContext is an interface to support dynamic dispatch.
type IShowDatabaseStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowDatabaseStmtContext differentiates from other interfaces.
	IsShowDatabaseStmtContext()
}

type ShowDatabaseStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowDatabaseStmtContext() *ShowDatabaseStmtContext {
	var p = new(ShowDatabaseStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showDatabaseStmt
	return p
}

func (*ShowDatabaseStmtContext) IsShowDatabaseStmtContext() {}

func NewShowDatabaseStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowDatabaseStmtContext {
	var p = new(ShowDatabaseStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showDatabaseStmt

	return p
}

func (s *ShowDatabaseStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowDatabaseStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowDatabaseStmtContext) T_DATASBAES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAES, 0)
}

func (s *ShowDatabaseStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowDatabaseStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowDatabaseStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowDatabaseStmt(s)
	}
}

func (s *ShowDatabaseStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowDatabaseStmt(s)
	}
}




func (p *SQLParser) ShowDatabaseStmt() (localctx IShowDatabaseStmtContext) {
	localctx = NewShowDatabaseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, SQLParserRULE_showDatabaseStmt)

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
		p.SetState(128)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(129)
		p.Match(SQLParserT_DATASBAES)
	}



	return localctx
}


// IShowNameSpacesStmtContext is an interface to support dynamic dispatch.
type IShowNameSpacesStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowNameSpacesStmtContext differentiates from other interfaces.
	IsShowNameSpacesStmtContext()
}

type ShowNameSpacesStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowNameSpacesStmtContext() *ShowNameSpacesStmtContext {
	var p = new(ShowNameSpacesStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showNameSpacesStmt
	return p
}

func (*ShowNameSpacesStmtContext) IsShowNameSpacesStmtContext() {}

func NewShowNameSpacesStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowNameSpacesStmtContext {
	var p = new(ShowNameSpacesStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showNameSpacesStmt

	return p
}

func (s *ShowNameSpacesStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowNameSpacesStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowNameSpacesStmtContext) T_NAMESPACES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NAMESPACES, 0)
}

func (s *ShowNameSpacesStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowNameSpacesStmtContext) T_NAMESPACE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NAMESPACE, 0)
}

func (s *ShowNameSpacesStmtContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *ShowNameSpacesStmtContext) Prefix() IPrefixContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrefixContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrefixContext)
}

func (s *ShowNameSpacesStmtContext) LimitClause() ILimitClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimitClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ShowNameSpacesStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowNameSpacesStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowNameSpacesStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowNameSpacesStmt(s)
	}
}

func (s *ShowNameSpacesStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowNameSpacesStmt(s)
	}
}




func (p *SQLParser) ShowNameSpacesStmt() (localctx IShowNameSpacesStmtContext) {
	localctx = NewShowNameSpacesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SQLParserRULE_showNameSpacesStmt)
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
		p.SetState(131)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(132)
		p.Match(SQLParserT_NAMESPACES)
	}
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(133)
			p.Match(SQLParserT_WHERE)
		}
		{
			p.SetState(134)
			p.Match(SQLParserT_NAMESPACE)
		}
		{
			p.SetState(135)
			p.Match(SQLParserT_EQUAL)
		}
		{
			p.SetState(136)
			p.Prefix()
		}

	}
	p.SetState(140)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(139)
			p.LimitClause()
		}

	}



	return localctx
}


// IShowMetricsStmtContext is an interface to support dynamic dispatch.
type IShowMetricsStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowMetricsStmtContext differentiates from other interfaces.
	IsShowMetricsStmtContext()
}

type ShowMetricsStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowMetricsStmtContext() *ShowMetricsStmtContext {
	var p = new(ShowMetricsStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showMetricsStmt
	return p
}

func (*ShowMetricsStmtContext) IsShowMetricsStmtContext() {}

func NewShowMetricsStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowMetricsStmtContext {
	var p = new(ShowMetricsStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showMetricsStmt

	return p
}

func (s *ShowMetricsStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowMetricsStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowMetricsStmtContext) T_METRICS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRICS, 0)
}

func (s *ShowMetricsStmtContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *ShowMetricsStmtContext) Namespace() INamespaceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INamespaceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
}

func (s *ShowMetricsStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowMetricsStmtContext) T_METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRIC, 0)
}

func (s *ShowMetricsStmtContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *ShowMetricsStmtContext) Prefix() IPrefixContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrefixContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrefixContext)
}

func (s *ShowMetricsStmtContext) LimitClause() ILimitClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILimitClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ShowMetricsStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMetricsStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowMetricsStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowMetricsStmt(s)
	}
}

func (s *ShowMetricsStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowMetricsStmt(s)
	}
}




func (p *SQLParser) ShowMetricsStmt() (localctx IShowMetricsStmtContext) {
	localctx = NewShowMetricsStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SQLParserRULE_showMetricsStmt)
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
		p.SetState(142)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(143)
		p.Match(SQLParserT_METRICS)
	}
	p.SetState(146)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(144)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(145)
			p.Namespace()
		}

	}
	p.SetState(152)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(148)
			p.Match(SQLParserT_WHERE)
		}
		{
			p.SetState(149)
			p.Match(SQLParserT_METRIC)
		}
		{
			p.SetState(150)
			p.Match(SQLParserT_EQUAL)
		}
		{
			p.SetState(151)
			p.Prefix()
		}

	}
	p.SetState(155)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(154)
			p.LimitClause()
		}

	}



	return localctx
}


// IShowFieldsStmtContext is an interface to support dynamic dispatch.
type IShowFieldsStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowFieldsStmtContext differentiates from other interfaces.
	IsShowFieldsStmtContext()
}

type ShowFieldsStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowFieldsStmtContext() *ShowFieldsStmtContext {
	var p = new(ShowFieldsStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showFieldsStmt
	return p
}

func (*ShowFieldsStmtContext) IsShowFieldsStmtContext() {}

func NewShowFieldsStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowFieldsStmtContext {
	var p = new(ShowFieldsStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showFieldsStmt

	return p
}

func (s *ShowFieldsStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowFieldsStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowFieldsStmtContext) T_FIELDS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIELDS, 0)
}

func (s *ShowFieldsStmtContext) FromClause() IFromClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFromClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFromClauseContext)
}

func (s *ShowFieldsStmtContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *ShowFieldsStmtContext) Namespace() INamespaceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INamespaceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
}

func (s *ShowFieldsStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowFieldsStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowFieldsStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowFieldsStmt(s)
	}
}

func (s *ShowFieldsStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowFieldsStmt(s)
	}
}




func (p *SQLParser) ShowFieldsStmt() (localctx IShowFieldsStmtContext) {
	localctx = NewShowFieldsStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SQLParserRULE_showFieldsStmt)
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
		p.SetState(157)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(158)
		p.Match(SQLParserT_FIELDS)
	}
	p.SetState(161)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(159)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(160)
			p.Namespace()
		}

	}
	{
		p.SetState(163)
		p.FromClause()
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

func (s *ShowTagKeysStmtContext) FromClause() IFromClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFromClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFromClauseContext)
}

func (s *ShowTagKeysStmtContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *ShowTagKeysStmtContext) Namespace() INamespaceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INamespaceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
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




func (p *SQLParser) ShowTagKeysStmt() (localctx IShowTagKeysStmtContext) {
	localctx = NewShowTagKeysStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SQLParserRULE_showTagKeysStmt)
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
		p.SetState(165)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(166)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(167)
		p.Match(SQLParserT_KEYS)
	}
	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(168)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(169)
			p.Namespace()
		}

	}
	{
		p.SetState(172)
		p.FromClause()
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

func (s *ShowTagValuesStmtContext) FromClause() IFromClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFromClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFromClauseContext)
}

func (s *ShowTagValuesStmtContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *ShowTagValuesStmtContext) T_KEY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEY, 0)
}

func (s *ShowTagValuesStmtContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *ShowTagValuesStmtContext) WithTagKey() IWithTagKeyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithTagKeyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWithTagKeyContext)
}

func (s *ShowTagValuesStmtContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *ShowTagValuesStmtContext) Namespace() INamespaceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INamespaceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
}

func (s *ShowTagValuesStmtContext) WhereClause() IWhereClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWhereClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWhereClauseContext)
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




func (p *SQLParser) ShowTagValuesStmt() (localctx IShowTagValuesStmtContext) {
	localctx = NewShowTagValuesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SQLParserRULE_showTagValuesStmt)
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
		p.SetState(174)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(175)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(176)
		p.Match(SQLParserT_VALUES)
	}
	p.SetState(179)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(177)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(178)
			p.Namespace()
		}

	}
	{
		p.SetState(181)
		p.FromClause()
	}
	{
		p.SetState(182)
		p.Match(SQLParserT_WITH)
	}
	{
		p.SetState(183)
		p.Match(SQLParserT_KEY)
	}
	{
		p.SetState(184)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(185)
		p.WithTagKey()
	}
	p.SetState(187)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(186)
			p.WhereClause()
		}

	}
	p.SetState(190)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(189)
			p.LimitClause()
		}

	}



	return localctx
}


// IPrefixContext is an interface to support dynamic dispatch.
type IPrefixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPrefixContext differentiates from other interfaces.
	IsPrefixContext()
}

type PrefixContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrefixContext() *PrefixContext {
	var p = new(PrefixContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_prefix
	return p
}

func (*PrefixContext) IsPrefixContext() {}

func NewPrefixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrefixContext {
	var p = new(PrefixContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_prefix

	return p
}

func (s *PrefixContext) GetParser() antlr.Parser { return s.parser }

func (s *PrefixContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *PrefixContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrefixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PrefixContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterPrefix(s)
	}
}

func (s *PrefixContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitPrefix(s)
	}
}




func (p *SQLParser) Prefix() (localctx IPrefixContext) {
	localctx = NewPrefixContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SQLParserRULE_prefix)

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
		p.SetState(192)
		p.Ident()
	}



	return localctx
}


// IWithTagKeyContext is an interface to support dynamic dispatch.
type IWithTagKeyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWithTagKeyContext differentiates from other interfaces.
	IsWithTagKeyContext()
}

type WithTagKeyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithTagKeyContext() *WithTagKeyContext {
	var p = new(WithTagKeyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_withTagKey
	return p
}

func (*WithTagKeyContext) IsWithTagKeyContext() {}

func NewWithTagKeyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithTagKeyContext {
	var p = new(WithTagKeyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_withTagKey

	return p
}

func (s *WithTagKeyContext) GetParser() antlr.Parser { return s.parser }

func (s *WithTagKeyContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *WithTagKeyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithTagKeyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WithTagKeyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWithTagKey(s)
	}
}

func (s *WithTagKeyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWithTagKey(s)
	}
}




func (p *SQLParser) WithTagKey() (localctx IWithTagKeyContext) {
	localctx = NewWithTagKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SQLParserRULE_withTagKey)

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
		p.Ident()
	}



	return localctx
}


// INamespaceContext is an interface to support dynamic dispatch.
type INamespaceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNamespaceContext differentiates from other interfaces.
	IsNamespaceContext()
}

type NamespaceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamespaceContext() *NamespaceContext {
	var p = new(NamespaceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_namespace
	return p
}

func (*NamespaceContext) IsNamespaceContext() {}

func NewNamespaceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamespaceContext {
	var p = new(NamespaceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_namespace

	return p
}

func (s *NamespaceContext) GetParser() antlr.Parser { return s.parser }

func (s *NamespaceContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *NamespaceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamespaceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *NamespaceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNamespace(s)
	}
}

func (s *NamespaceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNamespace(s)
	}
}




func (p *SQLParser) Namespace() (localctx INamespaceContext) {
	localctx = NewNamespaceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SQLParserRULE_namespace)

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
		p.SetState(196)
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

func (s *QueryStmtContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *QueryStmtContext) Namespace() INamespaceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INamespaceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
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
	p.EnterRule(localctx, 24, SQLParserRULE_queryStmt)
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
	p.SetState(199)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_EXPLAIN {
		{
			p.SetState(198)
			p.Match(SQLParserT_EXPLAIN)
		}

	}
	{
		p.SetState(201)
		p.SelectExpr()
	}
	p.SetState(204)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(202)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(203)
			p.Namespace()
		}

	}
	{
		p.SetState(206)
		p.FromClause()
	}
	p.SetState(208)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(207)
			p.WhereClause()
		}

	}
	p.SetState(211)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_GROUP {
		{
			p.SetState(210)
			p.GroupByClause()
		}

	}
	p.SetState(214)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ORDER {
		{
			p.SetState(213)
			p.OrderByClause()
		}

	}
	p.SetState(217)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(216)
			p.LimitClause()
		}

	}
	p.SetState(220)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WITH_VALUE {
		{
			p.SetState(219)
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
	p.EnterRule(localctx, 26, SQLParserRULE_selectExpr)

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
		p.SetState(222)
		p.Match(SQLParserT_SELECT)
	}
	{
		p.SetState(223)
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
	p.EnterRule(localctx, 28, SQLParserRULE_fields)
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
		p.SetState(225)
		p.Field()
	}
	p.SetState(230)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(226)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(227)
			p.Field()
		}


		p.SetState(232)
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
	p.EnterRule(localctx, 30, SQLParserRULE_field)
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
		p.SetState(233)
		p.fieldExpr(0)
	}
	p.SetState(235)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_AS {
		{
			p.SetState(234)
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
	p.EnterRule(localctx, 32, SQLParserRULE_alias)

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
		p.SetState(237)
		p.Match(SQLParserT_AS)
	}
	{
		p.SetState(238)
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
	p.EnterRule(localctx, 34, SQLParserRULE_fromClause)

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
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(241)
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
	p.EnterRule(localctx, 36, SQLParserRULE_whereClause)

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
		p.SetState(243)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(244)
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

func (s *ConditionExprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
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
	p.EnterRule(localctx, 38, SQLParserRULE_conditionExpr)
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

	p.SetState(256)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(246)
			p.tagFilterExpr(0)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(247)
			p.tagFilterExpr(0)
		}
		{
			p.SetState(248)
			p.Match(SQLParserT_AND)
		}
		{
			p.SetState(249)
			p.TimeRangeExpr()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(251)
			p.TimeRangeExpr()
		}
		p.SetState(254)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == SQLParserT_AND {
			{
				p.SetState(252)
				p.Match(SQLParserT_AND)
			}
			{
				p.SetState(253)
				p.tagFilterExpr(0)
			}

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
	_startState := 40
	p.EnterRecursionRule(localctx, 40, SQLParserRULE_tagFilterExpr, _p)
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
	p.SetState(286)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(259)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(260)
			p.tagFilterExpr(0)
		}
		{
			p.SetState(261)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(263)
			p.TagKey()
		}
		p.SetState(272)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_EQUAL:
			{
				p.SetState(264)
				p.Match(SQLParserT_EQUAL)
			}


		case SQLParserT_LIKE:
			{
				p.SetState(265)
				p.Match(SQLParserT_LIKE)
			}


		case SQLParserT_NOT:
			{
				p.SetState(266)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(267)
				p.Match(SQLParserT_LIKE)
			}


		case SQLParserT_REGEXP:
			{
				p.SetState(268)
				p.Match(SQLParserT_REGEXP)
			}


		case SQLParserT_NEQREGEXP:
			{
				p.SetState(269)
				p.Match(SQLParserT_NEQREGEXP)
			}


		case SQLParserT_NOTEQUAL:
			{
				p.SetState(270)
				p.Match(SQLParserT_NOTEQUAL)
			}


		case SQLParserT_NOTEQUAL2:
			{
				p.SetState(271)
				p.Match(SQLParserT_NOTEQUAL2)
			}



		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(274)
			p.TagValue()
		}


	case 3:
		{
			p.SetState(276)
			p.TagKey()
		}
		p.SetState(280)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_IN:
			{
				p.SetState(277)
				p.Match(SQLParserT_IN)
			}


		case SQLParserT_NOT:
			{
				p.SetState(278)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(279)
				p.Match(SQLParserT_IN)
			}



		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(282)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(283)
			p.TagValueList()
		}
		{
			p.SetState(284)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(293)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewTagFilterExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_tagFilterExpr)
			p.SetState(288)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(289)
				_la = p.GetTokenStream().LA(1)

				if !(_la == SQLParserT_AND || _la == SQLParserT_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(290)
				p.tagFilterExpr(2)
			}


		}
		p.SetState(295)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 42, SQLParserRULE_tagValueList)
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
		p.SetState(296)
		p.TagValue()
	}
	p.SetState(301)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(297)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(298)
			p.TagValue()
		}


		p.SetState(303)
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
	p.EnterRule(localctx, 44, SQLParserRULE_timeRangeExpr)

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
		p.SetState(304)
		p.TimeExpr()
	}
	p.SetState(307)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(305)
			p.Match(SQLParserT_AND)
		}
		{
			p.SetState(306)
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
	p.EnterRule(localctx, 46, SQLParserRULE_timeExpr)

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
		p.SetState(309)
		p.Match(SQLParserT_TIME)
	}
	{
		p.SetState(310)
		p.BinaryOperator()
	}
	p.SetState(313)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(311)
			p.NowExpr()
		}


	case 2:
		{
			p.SetState(312)
			p.Ident()
		}

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
	p.EnterRule(localctx, 48, SQLParserRULE_nowExpr)
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
		p.NowFunc()
	}
	p.SetState(317)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if ((((_la - 96)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 96))) & ((1 << (SQLParserT_ADD - 96)) | (1 << (SQLParserT_SUB - 96)) | (1 << (SQLParserL_INT - 96)))) != 0) {
		{
			p.SetState(316)
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
	p.EnterRule(localctx, 50, SQLParserRULE_nowFunc)
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
		p.SetState(319)
		p.Match(SQLParserT_NOW)
	}
	{
		p.SetState(320)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(322)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_DATASBAE) | (1 << SQLParserT_DATASBAES) | (1 << SQLParserT_NAMESPACE) | (1 << SQLParserT_NAMESPACES) | (1 << SQLParserT_NODE) | (1 << SQLParserT_METRICS) | (1 << SQLParserT_METRIC) | (1 << SQLParserT_FIELD) | (1 << SQLParserT_FIELDS) | (1 << SQLParserT_TAG) | (1 << SQLParserT_INFO) | (1 << SQLParserT_KEYS) | (1 << SQLParserT_KEY) | (1 << SQLParserT_WITH) | (1 << SQLParserT_VALUES))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)) | (1 << (SQLParserT_ASC - 32)) | (1 << (SQLParserT_DESC - 32)) | (1 << (SQLParserT_LIKE - 32)) | (1 << (SQLParserT_NOT - 32)) | (1 << (SQLParserT_BETWEEN - 32)) | (1 << (SQLParserT_IS - 32)) | (1 << (SQLParserT_GROUP - 32)) | (1 << (SQLParserT_HAVING - 32)) | (1 << (SQLParserT_BY - 32)) | (1 << (SQLParserT_FOR - 32)) | (1 << (SQLParserT_STATS - 32)) | (1 << (SQLParserT_TIME - 32)) | (1 << (SQLParserT_NOW - 32)) | (1 << (SQLParserT_IN - 32)) | (1 << (SQLParserT_LOG - 32)) | (1 << (SQLParserT_PROFILE - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)) | (1 << (SQLParserT_OPEN_P - 64)))) != 0) || ((((_la - 96)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 96))) & ((1 << (SQLParserT_ADD - 96)) | (1 << (SQLParserT_SUB - 96)) | (1 << (SQLParserL_ID - 96)) | (1 << (SQLParserL_INT - 96)) | (1 << (SQLParserL_DEC - 96)))) != 0) {
		{
			p.SetState(321)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(324)
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
	p.EnterRule(localctx, 52, SQLParserRULE_groupByClause)
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
		p.SetState(326)
		p.Match(SQLParserT_GROUP)
	}
	{
		p.SetState(327)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(328)
		p.GroupByKeys()
	}
	p.SetState(334)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_FILL {
		{
			p.SetState(329)
			p.Match(SQLParserT_FILL)
		}
		{
			p.SetState(330)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(331)
			p.FillOption()
		}
		{
			p.SetState(332)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.SetState(337)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_HAVING {
		{
			p.SetState(336)
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
	p.EnterRule(localctx, 54, SQLParserRULE_groupByKeys)
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
		p.SetState(339)
		p.GroupByKey()
	}
	p.SetState(344)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(340)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(341)
			p.GroupByKey()
		}


		p.SetState(346)
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
	p.EnterRule(localctx, 56, SQLParserRULE_groupByKey)

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

	p.SetState(353)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(347)
			p.Ident()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(348)
			p.Match(SQLParserT_TIME)
		}
		{
			p.SetState(349)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(350)
			p.DurationLit()
		}
		{
			p.SetState(351)
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
	p.EnterRule(localctx, 58, SQLParserRULE_fillOption)
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
		p.SetState(355)
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
	p.EnterRule(localctx, 60, SQLParserRULE_orderByClause)

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
		p.SetState(357)
		p.Match(SQLParserT_ORDER)
	}
	{
		p.SetState(358)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(359)
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
	p.EnterRule(localctx, 62, SQLParserRULE_sortField)
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
		p.SetState(361)
		p.fieldExpr(0)
	}
	p.SetState(365)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_ASC || _la == SQLParserT_DESC {
		{
			p.SetState(362)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_ASC || _la == SQLParserT_DESC) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}


		p.SetState(367)
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
	p.EnterRule(localctx, 64, SQLParserRULE_sortFields)
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
		p.SetState(368)
		p.SortField()
	}
	p.SetState(373)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(369)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(370)
			p.SortField()
		}


		p.SetState(375)
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
	p.EnterRule(localctx, 66, SQLParserRULE_havingClause)

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
		p.SetState(376)
		p.Match(SQLParserT_HAVING)
	}
	{
		p.SetState(377)
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
	_startState := 68
	p.EnterRecursionRule(localctx, 68, SQLParserRULE_boolExpr, _p)

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
	p.SetState(385)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 37, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(380)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(381)
			p.boolExpr(0)
		}
		{
			p.SetState(382)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(384)
			p.BoolExprAtom()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(393)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewBoolExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_boolExpr)
			p.SetState(387)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(388)
				p.BoolExprLogicalOp()
			}
			{
				p.SetState(389)
				p.boolExpr(3)
			}


		}
		p.SetState(395)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 70, SQLParserRULE_boolExprLogicalOp)
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
		p.SetState(396)
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
	p.EnterRule(localctx, 72, SQLParserRULE_boolExprAtom)

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
		p.SetState(398)
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
	p.EnterRule(localctx, 74, SQLParserRULE_binaryExpr)

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
		p.SetState(400)
		p.fieldExpr(0)
	}
	{
		p.SetState(401)
		p.BinaryOperator()
	}
	{
		p.SetState(402)
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
	p.EnterRule(localctx, 76, SQLParserRULE_binaryOperator)
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

	p.SetState(412)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_EQUAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(404)
			p.Match(SQLParserT_EQUAL)
		}


	case SQLParserT_NOTEQUAL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(405)
			p.Match(SQLParserT_NOTEQUAL)
		}


	case SQLParserT_NOTEQUAL2:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(406)
			p.Match(SQLParserT_NOTEQUAL2)
		}


	case SQLParserT_LESS:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(407)
			p.Match(SQLParserT_LESS)
		}


	case SQLParserT_LESSEQUAL:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(408)
			p.Match(SQLParserT_LESSEQUAL)
		}


	case SQLParserT_GREATER:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(409)
			p.Match(SQLParserT_GREATER)
		}


	case SQLParserT_GREATEREQUAL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(410)
			p.Match(SQLParserT_GREATEREQUAL)
		}


	case SQLParserT_LIKE, SQLParserT_REGEXP:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(411)
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
	_startState := 78
	p.EnterRecursionRule(localctx, 78, SQLParserRULE_fieldExpr, _p)

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
	p.SetState(422)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(415)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(416)
			p.fieldExpr(0)
		}
		{
			p.SetState(417)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(419)
			p.ExprFunc()
		}


	case 3:
		{
			p.SetState(420)
			p.ExprAtom()
		}


	case 4:
		{
			p.SetState(421)
			p.DurationLit()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(438)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(436)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext()) {
			case 1:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(424)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(425)
					p.Match(SQLParserT_MUL)
				}
				{
					p.SetState(426)
					p.fieldExpr(9)
				}


			case 2:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(427)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(428)
					p.Match(SQLParserT_DIV)
				}
				{
					p.SetState(429)
					p.fieldExpr(8)
				}


			case 3:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(430)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				}
				{
					p.SetState(431)
					p.Match(SQLParserT_ADD)
				}
				{
					p.SetState(432)
					p.fieldExpr(7)
				}


			case 4:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(433)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				{
					p.SetState(434)
					p.Match(SQLParserT_SUB)
				}
				{
					p.SetState(435)
					p.fieldExpr(6)
				}

			}

		}
		p.SetState(440)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 80, SQLParserRULE_durationLit)

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
		p.SetState(441)
		p.IntNumber()
	}
	{
		p.SetState(442)
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
	p.EnterRule(localctx, 82, SQLParserRULE_intervalItem)
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
		p.SetState(444)
		_la = p.GetTokenStream().LA(1)

		if !(((((_la - 71)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 71))) & ((1 << (SQLParserT_SECOND - 71)) | (1 << (SQLParserT_MINUTE - 71)) | (1 << (SQLParserT_HOUR - 71)) | (1 << (SQLParserT_DAY - 71)) | (1 << (SQLParserT_WEEK - 71)) | (1 << (SQLParserT_MONTH - 71)) | (1 << (SQLParserT_YEAR - 71)))) != 0)) {
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

func (s *ExprFuncContext) FuncName() IFuncNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFuncNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFuncNameContext)
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
	p.EnterRule(localctx, 84, SQLParserRULE_exprFunc)
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
		p.SetState(446)
		p.FuncName()
	}
	{
		p.SetState(447)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(449)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_DATASBAE) | (1 << SQLParserT_DATASBAES) | (1 << SQLParserT_NAMESPACE) | (1 << SQLParserT_NAMESPACES) | (1 << SQLParserT_NODE) | (1 << SQLParserT_METRICS) | (1 << SQLParserT_METRIC) | (1 << SQLParserT_FIELD) | (1 << SQLParserT_FIELDS) | (1 << SQLParserT_TAG) | (1 << SQLParserT_INFO) | (1 << SQLParserT_KEYS) | (1 << SQLParserT_KEY) | (1 << SQLParserT_WITH) | (1 << SQLParserT_VALUES))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)) | (1 << (SQLParserT_ASC - 32)) | (1 << (SQLParserT_DESC - 32)) | (1 << (SQLParserT_LIKE - 32)) | (1 << (SQLParserT_NOT - 32)) | (1 << (SQLParserT_BETWEEN - 32)) | (1 << (SQLParserT_IS - 32)) | (1 << (SQLParserT_GROUP - 32)) | (1 << (SQLParserT_HAVING - 32)) | (1 << (SQLParserT_BY - 32)) | (1 << (SQLParserT_FOR - 32)) | (1 << (SQLParserT_STATS - 32)) | (1 << (SQLParserT_TIME - 32)) | (1 << (SQLParserT_NOW - 32)) | (1 << (SQLParserT_IN - 32)) | (1 << (SQLParserT_LOG - 32)) | (1 << (SQLParserT_PROFILE - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)) | (1 << (SQLParserT_OPEN_P - 64)))) != 0) || ((((_la - 96)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 96))) & ((1 << (SQLParserT_ADD - 96)) | (1 << (SQLParserT_SUB - 96)) | (1 << (SQLParserL_ID - 96)) | (1 << (SQLParserL_INT - 96)) | (1 << (SQLParserL_DEC - 96)))) != 0) {
		{
			p.SetState(448)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(451)
		p.Match(SQLParserT_CLOSE_P)
	}



	return localctx
}


// IFuncNameContext is an interface to support dynamic dispatch.
type IFuncNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFuncNameContext differentiates from other interfaces.
	IsFuncNameContext()
}

type FuncNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuncNameContext() *FuncNameContext {
	var p = new(FuncNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_funcName
	return p
}

func (*FuncNameContext) IsFuncNameContext() {}

func NewFuncNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncNameContext {
	var p = new(FuncNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_funcName

	return p
}

func (s *FuncNameContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncNameContext) T_SUM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUM, 0)
}

func (s *FuncNameContext) T_MIN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MIN, 0)
}

func (s *FuncNameContext) T_MAX() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MAX, 0)
}

func (s *FuncNameContext) T_AVG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AVG, 0)
}

func (s *FuncNameContext) T_COUNT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_COUNT, 0)
}

func (s *FuncNameContext) T_STDDEV() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STDDEV, 0)
}

func (s *FuncNameContext) T_QUANTILE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUANTILE, 0)
}

func (s *FuncNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FuncNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFuncName(s)
	}
}

func (s *FuncNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFuncName(s)
	}
}




func (p *SQLParser) FuncName() (localctx IFuncNameContext) {
	localctx = NewFuncNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, SQLParserRULE_funcName)
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
		p.SetState(453)
		_la = p.GetTokenStream().LA(1)

		if !(((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)))) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
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
	p.EnterRule(localctx, 88, SQLParserRULE_exprFuncParams)
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
		p.SetState(455)
		p.FuncParam()
	}
	p.SetState(460)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(456)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(457)
			p.FuncParam()
		}


		p.SetState(462)
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
	p.EnterRule(localctx, 90, SQLParserRULE_funcParam)

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

	p.SetState(465)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(463)
			p.fieldExpr(0)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(464)
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
	p.EnterRule(localctx, 92, SQLParserRULE_exprAtom)

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

	p.SetState(473)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(467)
			p.Ident()
		}
		p.SetState(469)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(468)
				p.IdentFilter()
			}


		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(471)
			p.DecNumber()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(472)
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
	p.EnterRule(localctx, 94, SQLParserRULE_identFilter)

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
		p.SetState(475)
		p.Match(SQLParserT_OPEN_SB)
	}
	{
		p.SetState(476)
		p.tagFilterExpr(0)
	}
	{
		p.SetState(477)
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
	p.EnterRule(localctx, 96, SQLParserRULE_intNumber)
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
	p.SetState(480)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(479)
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
		p.SetState(482)
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
	p.EnterRule(localctx, 98, SQLParserRULE_decNumber)
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
	p.SetState(485)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(484)
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
		p.SetState(487)
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
	p.EnterRule(localctx, 100, SQLParserRULE_limitClause)

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
		p.SetState(489)
		p.Match(SQLParserT_LIMIT)
	}
	{
		p.SetState(490)
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
	p.EnterRule(localctx, 102, SQLParserRULE_metricName)

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
	p.EnterRule(localctx, 104, SQLParserRULE_tagKey)

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
		p.SetState(494)
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
	p.EnterRule(localctx, 106, SQLParserRULE_tagValue)

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
		p.SetState(496)
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
	p.EnterRule(localctx, 108, SQLParserRULE_ident)

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
	p.SetState(500)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserL_ID:
		{
			p.SetState(498)
			p.Match(SQLParserL_ID)
		}


	case SQLParserT_CREATE, SQLParserT_UPDATE, SQLParserT_SET, SQLParserT_DROP, SQLParserT_INTERVAL, SQLParserT_INTERVAL_NAME, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_META_TTL, SQLParserT_PAST_TTL, SQLParserT_FUTURE_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NAMESPACE, SQLParserT_NAMESPACES, SQLParserT_NODE, SQLParserT_METRICS, SQLParserT_METRIC, SQLParserT_FIELD, SQLParserT_FIELDS, SQLParserT_TAG, SQLParserT_INFO, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_VALUE, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_EXPLAIN, SQLParserT_WITH_VALUE, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_HAVING, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_NOW, SQLParserT_IN, SQLParserT_LOG, SQLParserT_PROFILE, SQLParserT_SUM, SQLParserT_MIN, SQLParserT_MAX, SQLParserT_COUNT, SQLParserT_AVG, SQLParserT_STDDEV, SQLParserT_QUANTILE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
		{
			p.SetState(499)
			p.NonReservedWords()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(509)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 52, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(502)
				p.Match(SQLParserT_DOT)
			}
			p.SetState(505)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case SQLParserL_ID:
				{
					p.SetState(503)
					p.Match(SQLParserL_ID)
				}


			case SQLParserT_CREATE, SQLParserT_UPDATE, SQLParserT_SET, SQLParserT_DROP, SQLParserT_INTERVAL, SQLParserT_INTERVAL_NAME, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_META_TTL, SQLParserT_PAST_TTL, SQLParserT_FUTURE_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NAMESPACE, SQLParserT_NAMESPACES, SQLParserT_NODE, SQLParserT_METRICS, SQLParserT_METRIC, SQLParserT_FIELD, SQLParserT_FIELDS, SQLParserT_TAG, SQLParserT_INFO, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_VALUE, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_EXPLAIN, SQLParserT_WITH_VALUE, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_HAVING, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_NOW, SQLParserT_IN, SQLParserT_LOG, SQLParserT_PROFILE, SQLParserT_SUM, SQLParserT_MIN, SQLParserT_MAX, SQLParserT_COUNT, SQLParserT_AVG, SQLParserT_STDDEV, SQLParserT_QUANTILE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
				{
					p.SetState(504)
					p.NonReservedWords()
				}



			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}


		}
		p.SetState(511)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 52, p.GetParserRuleContext())
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

func (s *NonReservedWordsContext) T_UPDATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_UPDATE, 0)
}

func (s *NonReservedWordsContext) T_SET() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SET, 0)
}

func (s *NonReservedWordsContext) T_DROP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DROP, 0)
}

func (s *NonReservedWordsContext) T_INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL, 0)
}

func (s *NonReservedWordsContext) T_INTERVAL_NAME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL_NAME, 0)
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

func (s *NonReservedWordsContext) T_META_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_META_TTL, 0)
}

func (s *NonReservedWordsContext) T_PAST_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PAST_TTL, 0)
}

func (s *NonReservedWordsContext) T_FUTURE_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FUTURE_TTL, 0)
}

func (s *NonReservedWordsContext) T_KILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KILL, 0)
}

func (s *NonReservedWordsContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *NonReservedWordsContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *NonReservedWordsContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *NonReservedWordsContext) T_DATASBAES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAES, 0)
}

func (s *NonReservedWordsContext) T_NAMESPACE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NAMESPACE, 0)
}

func (s *NonReservedWordsContext) T_NAMESPACES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NAMESPACES, 0)
}

func (s *NonReservedWordsContext) T_NODE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NODE, 0)
}

func (s *NonReservedWordsContext) T_METRICS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRICS, 0)
}

func (s *NonReservedWordsContext) T_METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRIC, 0)
}

func (s *NonReservedWordsContext) T_FIELD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIELD, 0)
}

func (s *NonReservedWordsContext) T_FIELDS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIELDS, 0)
}

func (s *NonReservedWordsContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *NonReservedWordsContext) T_INFO() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INFO, 0)
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

func (s *NonReservedWordsContext) T_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUE, 0)
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

func (s *NonReservedWordsContext) T_EXPLAIN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EXPLAIN, 0)
}

func (s *NonReservedWordsContext) T_WITH_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH_VALUE, 0)
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

func (s *NonReservedWordsContext) T_FILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FILL, 0)
}

func (s *NonReservedWordsContext) T_NULL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NULL, 0)
}

func (s *NonReservedWordsContext) T_PREVIOUS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PREVIOUS, 0)
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

func (s *NonReservedWordsContext) T_GROUP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GROUP, 0)
}

func (s *NonReservedWordsContext) T_HAVING() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HAVING, 0)
}

func (s *NonReservedWordsContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *NonReservedWordsContext) T_FOR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FOR, 0)
}

func (s *NonReservedWordsContext) T_STATS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATS, 0)
}

func (s *NonReservedWordsContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *NonReservedWordsContext) T_NOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOW, 0)
}

func (s *NonReservedWordsContext) T_IN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IN, 0)
}

func (s *NonReservedWordsContext) T_LOG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LOG, 0)
}

func (s *NonReservedWordsContext) T_PROFILE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PROFILE, 0)
}

func (s *NonReservedWordsContext) T_SUM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUM, 0)
}

func (s *NonReservedWordsContext) T_MIN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MIN, 0)
}

func (s *NonReservedWordsContext) T_MAX() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MAX, 0)
}

func (s *NonReservedWordsContext) T_COUNT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_COUNT, 0)
}

func (s *NonReservedWordsContext) T_AVG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AVG, 0)
}

func (s *NonReservedWordsContext) T_STDDEV() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STDDEV, 0)
}

func (s *NonReservedWordsContext) T_QUANTILE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUANTILE, 0)
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
	p.EnterRule(localctx, 110, SQLParserRULE_nonReservedWords)
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
		_la = p.GetTokenStream().LA(1)

		if !((((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_DATASBAE) | (1 << SQLParserT_DATASBAES) | (1 << SQLParserT_NAMESPACE) | (1 << SQLParserT_NAMESPACES) | (1 << SQLParserT_NODE) | (1 << SQLParserT_METRICS) | (1 << SQLParserT_METRIC) | (1 << SQLParserT_FIELD) | (1 << SQLParserT_FIELDS) | (1 << SQLParserT_TAG) | (1 << SQLParserT_INFO) | (1 << SQLParserT_KEYS) | (1 << SQLParserT_KEY) | (1 << SQLParserT_WITH) | (1 << SQLParserT_VALUES))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)) | (1 << (SQLParserT_ASC - 32)) | (1 << (SQLParserT_DESC - 32)) | (1 << (SQLParserT_LIKE - 32)) | (1 << (SQLParserT_NOT - 32)) | (1 << (SQLParserT_BETWEEN - 32)) | (1 << (SQLParserT_IS - 32)) | (1 << (SQLParserT_GROUP - 32)) | (1 << (SQLParserT_HAVING - 32)) | (1 << (SQLParserT_BY - 32)) | (1 << (SQLParserT_FOR - 32)) | (1 << (SQLParserT_STATS - 32)) | (1 << (SQLParserT_TIME - 32)) | (1 << (SQLParserT_NOW - 32)) | (1 << (SQLParserT_IN - 32)) | (1 << (SQLParserT_LOG - 32)) | (1 << (SQLParserT_PROFILE - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)))) != 0)) {
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
	case 20:
			var t *TagFilterExprContext = nil
			if localctx != nil { t = localctx.(*TagFilterExprContext) }
			return p.TagFilterExpr_Sempred(t, predIndex)

	case 34:
			var t *BoolExprContext = nil
			if localctx != nil { t = localctx.(*BoolExprContext) }
			return p.BoolExpr_Sempred(t, predIndex)

	case 39:
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

