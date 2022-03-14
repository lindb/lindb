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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 116, 612, 
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
	4, 66, 9, 66, 4, 67, 9, 67, 4, 68, 9, 68, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 
	5, 3, 154, 10, 3, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 
	3, 7, 3, 7, 3, 7, 3, 7, 3, 8, 3, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 10, 
	3, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 
	12, 3, 12, 5, 12, 189, 10, 12, 3, 12, 5, 12, 192, 10, 12, 3, 13, 3, 13, 
	3, 13, 3, 13, 5, 13, 198, 10, 13, 3, 13, 3, 13, 3, 13, 3, 13, 5, 13, 204, 
	10, 13, 3, 13, 5, 13, 207, 10, 13, 3, 14, 3, 14, 3, 14, 3, 14, 5, 14, 213, 
	10, 14, 3, 14, 3, 14, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 5, 15, 222, 10, 
	15, 3, 15, 3, 15, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 5, 16, 231, 10, 16, 
	3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 5, 16, 239, 10, 16, 3, 16, 5, 
	16, 242, 10, 16, 3, 17, 3, 17, 3, 18, 3, 18, 3, 19, 3, 19, 3, 20, 5, 20, 
	251, 10, 20, 3, 20, 3, 20, 3, 20, 5, 20, 256, 10, 20, 3, 20, 3, 20, 5, 
	20, 260, 10, 20, 3, 20, 5, 20, 263, 10, 20, 3, 20, 5, 20, 266, 10, 20, 
	3, 20, 5, 20, 269, 10, 20, 3, 20, 5, 20, 272, 10, 20, 3, 21, 3, 21, 3, 
	21, 3, 22, 3, 22, 3, 22, 7, 22, 280, 10, 22, 12, 22, 14, 22, 283, 11, 22, 
	3, 23, 3, 23, 5, 23, 287, 10, 23, 3, 24, 3, 24, 3, 24, 3, 25, 3, 25, 3, 
	25, 3, 26, 3, 26, 3, 26, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 
	3, 27, 5, 27, 306, 10, 27, 5, 27, 308, 10, 27, 3, 28, 3, 28, 3, 28, 3, 
	28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 
	5, 28, 324, 10, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 5, 28, 332, 
	10, 28, 3, 28, 3, 28, 3, 28, 3, 28, 5, 28, 338, 10, 28, 3, 28, 3, 28, 3, 
	28, 7, 28, 343, 10, 28, 12, 28, 14, 28, 346, 11, 28, 3, 29, 3, 29, 3, 29, 
	7, 29, 351, 10, 29, 12, 29, 14, 29, 354, 11, 29, 3, 30, 3, 30, 3, 30, 5, 
	30, 359, 10, 30, 3, 31, 3, 31, 3, 31, 3, 31, 5, 31, 365, 10, 31, 3, 32, 
	3, 32, 5, 32, 369, 10, 32, 3, 33, 3, 33, 3, 33, 5, 33, 374, 10, 33, 3, 
	33, 3, 33, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 5, 34, 
	386, 10, 34, 3, 34, 5, 34, 389, 10, 34, 3, 35, 3, 35, 3, 35, 7, 35, 394, 
	10, 35, 12, 35, 14, 35, 397, 11, 35, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36, 
	3, 36, 5, 36, 405, 10, 36, 3, 37, 3, 37, 3, 38, 3, 38, 3, 38, 3, 38, 3, 
	39, 3, 39, 7, 39, 415, 10, 39, 12, 39, 14, 39, 418, 11, 39, 3, 40, 3, 40, 
	3, 40, 7, 40, 423, 10, 40, 12, 40, 14, 40, 426, 11, 40, 3, 41, 3, 41, 3, 
	41, 3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 5, 42, 437, 10, 42, 3, 42, 
	3, 42, 3, 42, 3, 42, 7, 42, 443, 10, 42, 12, 42, 14, 42, 446, 11, 42, 3, 
	43, 3, 43, 3, 44, 3, 44, 3, 45, 3, 45, 3, 45, 3, 45, 3, 46, 3, 46, 3, 46, 
	3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 5, 46, 464, 10, 46, 3, 47, 3, 47, 3, 
	47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 5, 47, 474, 10, 47, 3, 47, 3, 47, 
	3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 7, 
	47, 488, 10, 47, 12, 47, 14, 47, 491, 11, 47, 3, 48, 3, 48, 3, 48, 3, 49, 
	3, 49, 3, 50, 3, 50, 3, 50, 5, 50, 501, 10, 50, 3, 50, 3, 50, 3, 51, 3, 
	51, 3, 52, 3, 52, 3, 52, 7, 52, 510, 10, 52, 12, 52, 14, 52, 513, 11, 52, 
	3, 53, 3, 53, 5, 53, 517, 10, 53, 3, 54, 3, 54, 5, 54, 521, 10, 54, 3, 
	54, 3, 54, 5, 54, 525, 10, 54, 3, 55, 3, 55, 3, 55, 3, 55, 3, 56, 3, 56, 
	3, 57, 3, 57, 3, 57, 3, 57, 7, 57, 537, 10, 57, 12, 57, 14, 57, 540, 11, 
	57, 3, 57, 3, 57, 3, 57, 3, 57, 5, 57, 546, 10, 57, 3, 58, 3, 58, 3, 58, 
	3, 58, 3, 59, 3, 59, 3, 59, 3, 59, 7, 59, 556, 10, 59, 12, 59, 14, 59, 
	559, 11, 59, 3, 59, 3, 59, 3, 59, 3, 59, 5, 59, 565, 10, 59, 3, 60, 3, 
	60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 5, 60, 575, 10, 60, 3, 61, 
	5, 61, 578, 10, 61, 3, 61, 3, 61, 3, 62, 5, 62, 583, 10, 62, 3, 62, 3, 
	62, 3, 63, 3, 63, 3, 63, 3, 64, 3, 64, 3, 65, 3, 65, 3, 66, 3, 66, 3, 67, 
	3, 67, 5, 67, 598, 10, 67, 3, 67, 3, 67, 3, 67, 5, 67, 603, 10, 67, 7, 
	67, 605, 10, 67, 12, 67, 14, 67, 608, 11, 67, 3, 68, 3, 68, 3, 68, 2, 5, 
	54, 82, 92, 69, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 
	32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 
	68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 
	104, 106, 108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 130, 132, 
	134, 2, 11, 3, 2, 26, 27, 3, 2, 55, 56, 4, 2, 58, 59, 115, 116, 3, 2, 61, 
	62, 4, 2, 63, 63, 100, 100, 3, 2, 84, 90, 3, 2, 77, 83, 3, 2, 109, 110, 
	3, 2, 8, 90, 2, 638, 2, 136, 3, 2, 2, 2, 4, 153, 3, 2, 2, 2, 6, 155, 3, 
	2, 2, 2, 8, 158, 3, 2, 2, 2, 10, 161, 3, 2, 2, 2, 12, 164, 3, 2, 2, 2, 
	14, 168, 3, 2, 2, 2, 16, 172, 3, 2, 2, 2, 18, 175, 3, 2, 2, 2, 20, 179, 
	3, 2, 2, 2, 22, 182, 3, 2, 2, 2, 24, 193, 3, 2, 2, 2, 26, 208, 3, 2, 2, 
	2, 28, 216, 3, 2, 2, 2, 30, 225, 3, 2, 2, 2, 32, 243, 3, 2, 2, 2, 34, 245, 
	3, 2, 2, 2, 36, 247, 3, 2, 2, 2, 38, 250, 3, 2, 2, 2, 40, 273, 3, 2, 2, 
	2, 42, 276, 3, 2, 2, 2, 44, 284, 3, 2, 2, 2, 46, 288, 3, 2, 2, 2, 48, 291, 
	3, 2, 2, 2, 50, 294, 3, 2, 2, 2, 52, 307, 3, 2, 2, 2, 54, 337, 3, 2, 2, 
	2, 56, 347, 3, 2, 2, 2, 58, 355, 3, 2, 2, 2, 60, 360, 3, 2, 2, 2, 62, 366, 
	3, 2, 2, 2, 64, 370, 3, 2, 2, 2, 66, 377, 3, 2, 2, 2, 68, 390, 3, 2, 2, 
	2, 70, 404, 3, 2, 2, 2, 72, 406, 3, 2, 2, 2, 74, 408, 3, 2, 2, 2, 76, 412, 
	3, 2, 2, 2, 78, 419, 3, 2, 2, 2, 80, 427, 3, 2, 2, 2, 82, 436, 3, 2, 2, 
	2, 84, 447, 3, 2, 2, 2, 86, 449, 3, 2, 2, 2, 88, 451, 3, 2, 2, 2, 90, 463, 
	3, 2, 2, 2, 92, 473, 3, 2, 2, 2, 94, 492, 3, 2, 2, 2, 96, 495, 3, 2, 2, 
	2, 98, 497, 3, 2, 2, 2, 100, 504, 3, 2, 2, 2, 102, 506, 3, 2, 2, 2, 104, 
	516, 3, 2, 2, 2, 106, 524, 3, 2, 2, 2, 108, 526, 3, 2, 2, 2, 110, 530, 
	3, 2, 2, 2, 112, 545, 3, 2, 2, 2, 114, 547, 3, 2, 2, 2, 116, 564, 3, 2, 
	2, 2, 118, 574, 3, 2, 2, 2, 120, 577, 3, 2, 2, 2, 122, 582, 3, 2, 2, 2, 
	124, 586, 3, 2, 2, 2, 126, 589, 3, 2, 2, 2, 128, 591, 3, 2, 2, 2, 130, 
	593, 3, 2, 2, 2, 132, 597, 3, 2, 2, 2, 134, 609, 3, 2, 2, 2, 136, 137, 
	5, 4, 3, 2, 137, 138, 7, 2, 2, 3, 138, 3, 3, 2, 2, 2, 139, 154, 5, 8, 5, 
	2, 140, 154, 5, 10, 6, 2, 141, 154, 5, 12, 7, 2, 142, 154, 5, 14, 8, 2, 
	143, 154, 5, 16, 9, 2, 144, 154, 5, 20, 11, 2, 145, 154, 5, 6, 4, 2, 146, 
	154, 5, 22, 12, 2, 147, 154, 5, 24, 13, 2, 148, 154, 5, 26, 14, 2, 149, 
	154, 5, 28, 15, 2, 150, 154, 5, 30, 16, 2, 151, 154, 5, 38, 20, 2, 152, 
	154, 5, 18, 10, 2, 153, 139, 3, 2, 2, 2, 153, 140, 3, 2, 2, 2, 153, 141, 
	3, 2, 2, 2, 153, 142, 3, 2, 2, 2, 153, 143, 3, 2, 2, 2, 153, 144, 3, 2, 
	2, 2, 153, 145, 3, 2, 2, 2, 153, 146, 3, 2, 2, 2, 153, 147, 3, 2, 2, 2, 
	153, 148, 3, 2, 2, 2, 153, 149, 3, 2, 2, 2, 153, 150, 3, 2, 2, 2, 153, 
	151, 3, 2, 2, 2, 153, 152, 3, 2, 2, 2, 154, 5, 3, 2, 2, 2, 155, 156, 7, 
	23, 2, 2, 156, 157, 5, 132, 67, 2, 157, 7, 3, 2, 2, 2, 158, 159, 7, 22, 
	2, 2, 159, 160, 7, 24, 2, 2, 160, 9, 3, 2, 2, 2, 161, 162, 7, 22, 2, 2, 
	162, 163, 7, 25, 2, 2, 163, 11, 3, 2, 2, 2, 164, 165, 7, 22, 2, 2, 165, 
	166, 7, 28, 2, 2, 166, 167, 9, 2, 2, 2, 167, 13, 3, 2, 2, 2, 168, 169, 
	7, 8, 2, 2, 169, 170, 7, 26, 2, 2, 170, 171, 5, 110, 56, 2, 171, 15, 3, 
	2, 2, 2, 172, 173, 7, 22, 2, 2, 173, 174, 7, 29, 2, 2, 174, 17, 3, 2, 2, 
	2, 175, 176, 7, 8, 2, 2, 176, 177, 7, 30, 2, 2, 177, 178, 5, 110, 56, 2, 
	178, 19, 3, 2, 2, 2, 179, 180, 7, 22, 2, 2, 180, 181, 7, 31, 2, 2, 181, 
	21, 3, 2, 2, 2, 182, 183, 7, 22, 2, 2, 183, 188, 7, 33, 2, 2, 184, 185, 
	7, 47, 2, 2, 185, 186, 7, 32, 2, 2, 186, 187, 7, 93, 2, 2, 187, 189, 5, 
	32, 17, 2, 188, 184, 3, 2, 2, 2, 188, 189, 3, 2, 2, 2, 189, 191, 3, 2, 
	2, 2, 190, 192, 5, 124, 63, 2, 191, 190, 3, 2, 2, 2, 191, 192, 3, 2, 2, 
	2, 192, 23, 3, 2, 2, 2, 193, 194, 7, 22, 2, 2, 194, 197, 7, 35, 2, 2, 195, 
	196, 7, 21, 2, 2, 196, 198, 5, 36, 19, 2, 197, 195, 3, 2, 2, 2, 197, 198, 
	3, 2, 2, 2, 198, 203, 3, 2, 2, 2, 199, 200, 7, 47, 2, 2, 200, 201, 7, 36, 
	2, 2, 201, 202, 7, 93, 2, 2, 202, 204, 5, 32, 17, 2, 203, 199, 3, 2, 2, 
	2, 203, 204, 3, 2, 2, 2, 204, 206, 3, 2, 2, 2, 205, 207, 5, 124, 63, 2, 
	206, 205, 3, 2, 2, 2, 206, 207, 3, 2, 2, 2, 207, 25, 3, 2, 2, 2, 208, 209, 
	7, 22, 2, 2, 209, 212, 7, 38, 2, 2, 210, 211, 7, 21, 2, 2, 211, 213, 5, 
	36, 19, 2, 212, 210, 3, 2, 2, 2, 212, 213, 3, 2, 2, 2, 213, 214, 3, 2, 
	2, 2, 214, 215, 5, 48, 25, 2, 215, 27, 3, 2, 2, 2, 216, 217, 7, 22, 2, 
	2, 217, 218, 7, 39, 2, 2, 218, 221, 7, 41, 2, 2, 219, 220, 7, 21, 2, 2, 
	220, 222, 5, 36, 19, 2, 221, 219, 3, 2, 2, 2, 221, 222, 3, 2, 2, 2, 222, 
	223, 3, 2, 2, 2, 223, 224, 5, 48, 25, 2, 224, 29, 3, 2, 2, 2, 225, 226, 
	7, 22, 2, 2, 226, 227, 7, 39, 2, 2, 227, 230, 7, 44, 2, 2, 228, 229, 7, 
	21, 2, 2, 229, 231, 5, 36, 19, 2, 230, 228, 3, 2, 2, 2, 230, 231, 3, 2, 
	2, 2, 231, 232, 3, 2, 2, 2, 232, 233, 5, 48, 25, 2, 233, 234, 7, 43, 2, 
	2, 234, 235, 7, 42, 2, 2, 235, 236, 7, 93, 2, 2, 236, 238, 5, 34, 18, 2, 
	237, 239, 5, 50, 26, 2, 238, 237, 3, 2, 2, 2, 238, 239, 3, 2, 2, 2, 239, 
	241, 3, 2, 2, 2, 240, 242, 5, 124, 63, 2, 241, 240, 3, 2, 2, 2, 241, 242, 
	3, 2, 2, 2, 242, 31, 3, 2, 2, 2, 243, 244, 5, 132, 67, 2, 244, 33, 3, 2, 
	2, 2, 245, 246, 5, 132, 67, 2, 246, 35, 3, 2, 2, 2, 247, 248, 5, 132, 67, 
	2, 248, 37, 3, 2, 2, 2, 249, 251, 7, 51, 2, 2, 250, 249, 3, 2, 2, 2, 250, 
	251, 3, 2, 2, 2, 251, 252, 3, 2, 2, 2, 252, 255, 5, 40, 21, 2, 253, 254, 
	7, 21, 2, 2, 254, 256, 5, 36, 19, 2, 255, 253, 3, 2, 2, 2, 255, 256, 3, 
	2, 2, 2, 256, 257, 3, 2, 2, 2, 257, 259, 5, 48, 25, 2, 258, 260, 5, 50, 
	26, 2, 259, 258, 3, 2, 2, 2, 259, 260, 3, 2, 2, 2, 260, 262, 3, 2, 2, 2, 
	261, 263, 5, 66, 34, 2, 262, 261, 3, 2, 2, 2, 262, 263, 3, 2, 2, 2, 263, 
	265, 3, 2, 2, 2, 264, 266, 5, 74, 38, 2, 265, 264, 3, 2, 2, 2, 265, 266, 
	3, 2, 2, 2, 266, 268, 3, 2, 2, 2, 267, 269, 5, 124, 63, 2, 268, 267, 3, 
	2, 2, 2, 268, 269, 3, 2, 2, 2, 269, 271, 3, 2, 2, 2, 270, 272, 7, 52, 2, 
	2, 271, 270, 3, 2, 2, 2, 271, 272, 3, 2, 2, 2, 272, 39, 3, 2, 2, 2, 273, 
	274, 7, 53, 2, 2, 274, 275, 5, 42, 22, 2, 275, 41, 3, 2, 2, 2, 276, 281, 
	5, 44, 23, 2, 277, 278, 7, 102, 2, 2, 278, 280, 5, 44, 23, 2, 279, 277, 
	3, 2, 2, 2, 280, 283, 3, 2, 2, 2, 281, 279, 3, 2, 2, 2, 281, 282, 3, 2, 
	2, 2, 282, 43, 3, 2, 2, 2, 283, 281, 3, 2, 2, 2, 284, 286, 5, 92, 47, 2, 
	285, 287, 5, 46, 24, 2, 286, 285, 3, 2, 2, 2, 286, 287, 3, 2, 2, 2, 287, 
	45, 3, 2, 2, 2, 288, 289, 7, 54, 2, 2, 289, 290, 5, 132, 67, 2, 290, 47, 
	3, 2, 2, 2, 291, 292, 7, 46, 2, 2, 292, 293, 5, 126, 64, 2, 293, 49, 3, 
	2, 2, 2, 294, 295, 7, 47, 2, 2, 295, 296, 5, 52, 27, 2, 296, 51, 3, 2, 
	2, 2, 297, 308, 5, 54, 28, 2, 298, 299, 5, 54, 28, 2, 299, 300, 7, 55, 
	2, 2, 300, 301, 5, 58, 30, 2, 301, 308, 3, 2, 2, 2, 302, 305, 5, 58, 30, 
	2, 303, 304, 7, 55, 2, 2, 304, 306, 5, 54, 28, 2, 305, 303, 3, 2, 2, 2, 
	305, 306, 3, 2, 2, 2, 306, 308, 3, 2, 2, 2, 307, 297, 3, 2, 2, 2, 307, 
	298, 3, 2, 2, 2, 307, 302, 3, 2, 2, 2, 308, 53, 3, 2, 2, 2, 309, 310, 8, 
	28, 1, 2, 310, 311, 7, 107, 2, 2, 311, 312, 5, 54, 28, 2, 312, 313, 7, 
	108, 2, 2, 313, 338, 3, 2, 2, 2, 314, 323, 5, 128, 65, 2, 315, 324, 7, 
	93, 2, 2, 316, 324, 7, 63, 2, 2, 317, 318, 7, 64, 2, 2, 318, 324, 7, 63, 
	2, 2, 319, 324, 7, 100, 2, 2, 320, 324, 7, 101, 2, 2, 321, 324, 7, 94, 
	2, 2, 322, 324, 7, 95, 2, 2, 323, 315, 3, 2, 2, 2, 323, 316, 3, 2, 2, 2, 
	323, 317, 3, 2, 2, 2, 323, 319, 3, 2, 2, 2, 323, 320, 3, 2, 2, 2, 323, 
	321, 3, 2, 2, 2, 323, 322, 3, 2, 2, 2, 324, 325, 3, 2, 2, 2, 325, 326, 
	5, 130, 66, 2, 326, 338, 3, 2, 2, 2, 327, 331, 5, 128, 65, 2, 328, 332, 
	7, 74, 2, 2, 329, 330, 7, 64, 2, 2, 330, 332, 7, 74, 2, 2, 331, 328, 3, 
	2, 2, 2, 331, 329, 3, 2, 2, 2, 332, 333, 3, 2, 2, 2, 333, 334, 7, 107, 
	2, 2, 334, 335, 5, 56, 29, 2, 335, 336, 7, 108, 2, 2, 336, 338, 3, 2, 2, 
	2, 337, 309, 3, 2, 2, 2, 337, 314, 3, 2, 2, 2, 337, 327, 3, 2, 2, 2, 338, 
	344, 3, 2, 2, 2, 339, 340, 12, 3, 2, 2, 340, 341, 9, 3, 2, 2, 341, 343, 
	5, 54, 28, 4, 342, 339, 3, 2, 2, 2, 343, 346, 3, 2, 2, 2, 344, 342, 3, 
	2, 2, 2, 344, 345, 3, 2, 2, 2, 345, 55, 3, 2, 2, 2, 346, 344, 3, 2, 2, 
	2, 347, 352, 5, 130, 66, 2, 348, 349, 7, 102, 2, 2, 349, 351, 5, 130, 66, 
	2, 350, 348, 3, 2, 2, 2, 351, 354, 3, 2, 2, 2, 352, 350, 3, 2, 2, 2, 352, 
	353, 3, 2, 2, 2, 353, 57, 3, 2, 2, 2, 354, 352, 3, 2, 2, 2, 355, 358, 5, 
	60, 31, 2, 356, 357, 7, 55, 2, 2, 357, 359, 5, 60, 31, 2, 358, 356, 3, 
	2, 2, 2, 358, 359, 3, 2, 2, 2, 359, 59, 3, 2, 2, 2, 360, 361, 7, 72, 2, 
	2, 361, 364, 5, 90, 46, 2, 362, 365, 5, 62, 32, 2, 363, 365, 5, 132, 67, 
	2, 364, 362, 3, 2, 2, 2, 364, 363, 3, 2, 2, 2, 365, 61, 3, 2, 2, 2, 366, 
	368, 5, 64, 33, 2, 367, 369, 5, 94, 48, 2, 368, 367, 3, 2, 2, 2, 368, 369, 
	3, 2, 2, 2, 369, 63, 3, 2, 2, 2, 370, 371, 7, 73, 2, 2, 371, 373, 7, 107, 
	2, 2, 372, 374, 5, 102, 52, 2, 373, 372, 3, 2, 2, 2, 373, 374, 3, 2, 2, 
	2, 374, 375, 3, 2, 2, 2, 375, 376, 7, 108, 2, 2, 376, 65, 3, 2, 2, 2, 377, 
	378, 7, 67, 2, 2, 378, 379, 7, 69, 2, 2, 379, 385, 5, 68, 35, 2, 380, 381, 
	7, 57, 2, 2, 381, 382, 7, 107, 2, 2, 382, 383, 5, 72, 37, 2, 383, 384, 
	7, 108, 2, 2, 384, 386, 3, 2, 2, 2, 385, 380, 3, 2, 2, 2, 385, 386, 3, 
	2, 2, 2, 386, 388, 3, 2, 2, 2, 387, 389, 5, 80, 41, 2, 388, 387, 3, 2, 
	2, 2, 388, 389, 3, 2, 2, 2, 389, 67, 3, 2, 2, 2, 390, 395, 5, 70, 36, 2, 
	391, 392, 7, 102, 2, 2, 392, 394, 5, 70, 36, 2, 393, 391, 3, 2, 2, 2, 394, 
	397, 3, 2, 2, 2, 395, 393, 3, 2, 2, 2, 395, 396, 3, 2, 2, 2, 396, 69, 3, 
	2, 2, 2, 397, 395, 3, 2, 2, 2, 398, 405, 5, 132, 67, 2, 399, 400, 7, 72, 
	2, 2, 400, 401, 7, 107, 2, 2, 401, 402, 5, 94, 48, 2, 402, 403, 7, 108, 
	2, 2, 403, 405, 3, 2, 2, 2, 404, 398, 3, 2, 2, 2, 404, 399, 3, 2, 2, 2, 
	405, 71, 3, 2, 2, 2, 406, 407, 9, 4, 2, 2, 407, 73, 3, 2, 2, 2, 408, 409, 
	7, 60, 2, 2, 409, 410, 7, 69, 2, 2, 410, 411, 5, 78, 40, 2, 411, 75, 3, 
	2, 2, 2, 412, 416, 5, 92, 47, 2, 413, 415, 9, 5, 2, 2, 414, 413, 3, 2, 
	2, 2, 415, 418, 3, 2, 2, 2, 416, 414, 3, 2, 2, 2, 416, 417, 3, 2, 2, 2, 
	417, 77, 3, 2, 2, 2, 418, 416, 3, 2, 2, 2, 419, 424, 5, 76, 39, 2, 420, 
	421, 7, 102, 2, 2, 421, 423, 5, 76, 39, 2, 422, 420, 3, 2, 2, 2, 423, 426, 
	3, 2, 2, 2, 424, 422, 3, 2, 2, 2, 424, 425, 3, 2, 2, 2, 425, 79, 3, 2, 
	2, 2, 426, 424, 3, 2, 2, 2, 427, 428, 7, 68, 2, 2, 428, 429, 5, 82, 42, 
	2, 429, 81, 3, 2, 2, 2, 430, 431, 8, 42, 1, 2, 431, 432, 7, 107, 2, 2, 
	432, 433, 5, 82, 42, 2, 433, 434, 7, 108, 2, 2, 434, 437, 3, 2, 2, 2, 435, 
	437, 5, 86, 44, 2, 436, 430, 3, 2, 2, 2, 436, 435, 3, 2, 2, 2, 437, 444, 
	3, 2, 2, 2, 438, 439, 12, 4, 2, 2, 439, 440, 5, 84, 43, 2, 440, 441, 5, 
	82, 42, 5, 441, 443, 3, 2, 2, 2, 442, 438, 3, 2, 2, 2, 443, 446, 3, 2, 
	2, 2, 444, 442, 3, 2, 2, 2, 444, 445, 3, 2, 2, 2, 445, 83, 3, 2, 2, 2, 
	446, 444, 3, 2, 2, 2, 447, 448, 9, 3, 2, 2, 448, 85, 3, 2, 2, 2, 449, 450, 
	5, 88, 45, 2, 450, 87, 3, 2, 2, 2, 451, 452, 5, 92, 47, 2, 452, 453, 5, 
	90, 46, 2, 453, 454, 5, 92, 47, 2, 454, 89, 3, 2, 2, 2, 455, 464, 7, 93, 
	2, 2, 456, 464, 7, 94, 2, 2, 457, 464, 7, 95, 2, 2, 458, 464, 7, 98, 2, 
	2, 459, 464, 7, 99, 2, 2, 460, 464, 7, 96, 2, 2, 461, 464, 7, 97, 2, 2, 
	462, 464, 9, 6, 2, 2, 463, 455, 3, 2, 2, 2, 463, 456, 3, 2, 2, 2, 463, 
	457, 3, 2, 2, 2, 463, 458, 3, 2, 2, 2, 463, 459, 3, 2, 2, 2, 463, 460, 
	3, 2, 2, 2, 463, 461, 3, 2, 2, 2, 463, 462, 3, 2, 2, 2, 464, 91, 3, 2, 
	2, 2, 465, 466, 8, 47, 1, 2, 466, 467, 7, 107, 2, 2, 467, 468, 5, 92, 47, 
	2, 468, 469, 7, 108, 2, 2, 469, 474, 3, 2, 2, 2, 470, 474, 5, 98, 50, 2, 
	471, 474, 5, 106, 54, 2, 472, 474, 5, 94, 48, 2, 473, 465, 3, 2, 2, 2, 
	473, 470, 3, 2, 2, 2, 473, 471, 3, 2, 2, 2, 473, 472, 3, 2, 2, 2, 474, 
	489, 3, 2, 2, 2, 475, 476, 12, 10, 2, 2, 476, 477, 7, 112, 2, 2, 477, 488, 
	5, 92, 47, 11, 478, 479, 12, 9, 2, 2, 479, 480, 7, 111, 2, 2, 480, 488, 
	5, 92, 47, 10, 481, 482, 12, 8, 2, 2, 482, 483, 7, 109, 2, 2, 483, 488, 
	5, 92, 47, 9, 484, 485, 12, 7, 2, 2, 485, 486, 7, 110, 2, 2, 486, 488, 
	5, 92, 47, 8, 487, 475, 3, 2, 2, 2, 487, 478, 3, 2, 2, 2, 487, 481, 3, 
	2, 2, 2, 487, 484, 3, 2, 2, 2, 488, 491, 3, 2, 2, 2, 489, 487, 3, 2, 2, 
	2, 489, 490, 3, 2, 2, 2, 490, 93, 3, 2, 2, 2, 491, 489, 3, 2, 2, 2, 492, 
	493, 5, 120, 61, 2, 493, 494, 5, 96, 49, 2, 494, 95, 3, 2, 2, 2, 495, 496, 
	9, 7, 2, 2, 496, 97, 3, 2, 2, 2, 497, 498, 5, 100, 51, 2, 498, 500, 7, 
	107, 2, 2, 499, 501, 5, 102, 52, 2, 500, 499, 3, 2, 2, 2, 500, 501, 3, 
	2, 2, 2, 501, 502, 3, 2, 2, 2, 502, 503, 7, 108, 2, 2, 503, 99, 3, 2, 2, 
	2, 504, 505, 9, 8, 2, 2, 505, 101, 3, 2, 2, 2, 506, 511, 5, 104, 53, 2, 
	507, 508, 7, 102, 2, 2, 508, 510, 5, 104, 53, 2, 509, 507, 3, 2, 2, 2, 
	510, 513, 3, 2, 2, 2, 511, 509, 3, 2, 2, 2, 511, 512, 3, 2, 2, 2, 512, 
	103, 3, 2, 2, 2, 513, 511, 3, 2, 2, 2, 514, 517, 5, 92, 47, 2, 515, 517, 
	5, 54, 28, 2, 516, 514, 3, 2, 2, 2, 516, 515, 3, 2, 2, 2, 517, 105, 3, 
	2, 2, 2, 518, 520, 5, 132, 67, 2, 519, 521, 5, 108, 55, 2, 520, 519, 3, 
	2, 2, 2, 520, 521, 3, 2, 2, 2, 521, 525, 3, 2, 2, 2, 522, 525, 5, 122, 
	62, 2, 523, 525, 5, 120, 61, 2, 524, 518, 3, 2, 2, 2, 524, 522, 3, 2, 2, 
	2, 524, 523, 3, 2, 2, 2, 525, 107, 3, 2, 2, 2, 526, 527, 7, 105, 2, 2, 
	527, 528, 5, 54, 28, 2, 528, 529, 7, 106, 2, 2, 529, 109, 3, 2, 2, 2, 530, 
	531, 5, 118, 60, 2, 531, 111, 3, 2, 2, 2, 532, 533, 7, 103, 2, 2, 533, 
	538, 5, 114, 58, 2, 534, 535, 7, 102, 2, 2, 535, 537, 5, 114, 58, 2, 536, 
	534, 3, 2, 2, 2, 537, 540, 3, 2, 2, 2, 538, 536, 3, 2, 2, 2, 538, 539, 
	3, 2, 2, 2, 539, 541, 3, 2, 2, 2, 540, 538, 3, 2, 2, 2, 541, 542, 7, 104, 
	2, 2, 542, 546, 3, 2, 2, 2, 543, 544, 7, 103, 2, 2, 544, 546, 7, 104, 2, 
	2, 545, 532, 3, 2, 2, 2, 545, 543, 3, 2, 2, 2, 546, 113, 3, 2, 2, 2, 547, 
	548, 7, 6, 2, 2, 548, 549, 7, 92, 2, 2, 549, 550, 5, 118, 60, 2, 550, 115, 
	3, 2, 2, 2, 551, 552, 7, 105, 2, 2, 552, 557, 5, 118, 60, 2, 553, 554, 
	7, 102, 2, 2, 554, 556, 5, 118, 60, 2, 555, 553, 3, 2, 2, 2, 556, 559, 
	3, 2, 2, 2, 557, 555, 3, 2, 2, 2, 557, 558, 3, 2, 2, 2, 558, 560, 3, 2, 
	2, 2, 559, 557, 3, 2, 2, 2, 560, 561, 7, 106, 2, 2, 561, 565, 3, 2, 2, 
	2, 562, 563, 7, 105, 2, 2, 563, 565, 7, 106, 2, 2, 564, 551, 3, 2, 2, 2, 
	564, 562, 3, 2, 2, 2, 565, 117, 3, 2, 2, 2, 566, 575, 7, 6, 2, 2, 567, 
	575, 5, 120, 61, 2, 568, 575, 5, 122, 62, 2, 569, 575, 5, 112, 57, 2, 570, 
	575, 5, 116, 59, 2, 571, 575, 7, 3, 2, 2, 572, 575, 7, 4, 2, 2, 573, 575, 
	7, 5, 2, 2, 574, 566, 3, 2, 2, 2, 574, 567, 3, 2, 2, 2, 574, 568, 3, 2, 
	2, 2, 574, 569, 3, 2, 2, 2, 574, 570, 3, 2, 2, 2, 574, 571, 3, 2, 2, 2, 
	574, 572, 3, 2, 2, 2, 574, 573, 3, 2, 2, 2, 575, 119, 3, 2, 2, 2, 576, 
	578, 9, 9, 2, 2, 577, 576, 3, 2, 2, 2, 577, 578, 3, 2, 2, 2, 578, 579, 
	3, 2, 2, 2, 579, 580, 7, 115, 2, 2, 580, 121, 3, 2, 2, 2, 581, 583, 9, 
	9, 2, 2, 582, 581, 3, 2, 2, 2, 582, 583, 3, 2, 2, 2, 583, 584, 3, 2, 2, 
	2, 584, 585, 7, 116, 2, 2, 585, 123, 3, 2, 2, 2, 586, 587, 7, 48, 2, 2, 
	587, 588, 7, 115, 2, 2, 588, 125, 3, 2, 2, 2, 589, 590, 5, 132, 67, 2, 
	590, 127, 3, 2, 2, 2, 591, 592, 5, 132, 67, 2, 592, 129, 3, 2, 2, 2, 593, 
	594, 5, 132, 67, 2, 594, 131, 3, 2, 2, 2, 595, 598, 7, 114, 2, 2, 596, 
	598, 5, 134, 68, 2, 597, 595, 3, 2, 2, 2, 597, 596, 3, 2, 2, 2, 598, 606, 
	3, 2, 2, 2, 599, 602, 7, 91, 2, 2, 600, 603, 7, 114, 2, 2, 601, 603, 5, 
	134, 68, 2, 602, 600, 3, 2, 2, 2, 602, 601, 3, 2, 2, 2, 603, 605, 3, 2, 
	2, 2, 604, 599, 3, 2, 2, 2, 605, 608, 3, 2, 2, 2, 606, 604, 3, 2, 2, 2, 
	606, 607, 3, 2, 2, 2, 607, 133, 3, 2, 2, 2, 608, 606, 3, 2, 2, 2, 609, 
	610, 9, 10, 2, 2, 610, 135, 3, 2, 2, 2, 60, 153, 188, 191, 197, 203, 206, 
	212, 221, 230, 238, 241, 250, 255, 259, 262, 265, 268, 271, 281, 286, 305, 
	307, 323, 331, 337, 344, 352, 358, 364, 368, 373, 385, 388, 395, 404, 416, 
	424, 436, 444, 463, 473, 487, 489, 500, 511, 516, 520, 524, 538, 545, 557, 
	564, 574, 577, 582, 597, 602, 606,
}
var literalNames = []string{
	"", "'true'", "'false'", "'null'", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "'m'", 
	"", "", "", "'M'", "", "'.'", "':'", "'='", "'<>'", "'!='", "'>'", "'>='", 
	"'<'", "'<='", "'=~'", "'!~'", "','", "'{'", "'}'", "'['", "']'", "'('", 
	"')'", "'+'", "'-'", "'/'", "'*'", "'%'",
}
var symbolicNames = []string{
	"", "", "", "", "STRING", "WS", "T_CREATE", "T_UPDATE", "T_SET", "T_DROP", 
	"T_INTERVAL", "T_INTERVAL_NAME", "T_SHARD", "T_REPLICATION", "T_TTL", "T_META_TTL", 
	"T_PAST_TTL", "T_FUTURE_TTL", "T_KILL", "T_ON", "T_SHOW", "T_USE", "T_MASTER", 
	"T_STORAGES", "T_STORAGE", "T_BROKER", "T_ALIVE", "T_SCHEMAS", "T_DATASBAE", 
	"T_DATASBAES", "T_NAMESPACE", "T_NAMESPACES", "T_NODE", "T_METRICS", "T_METRIC", 
	"T_FIELD", "T_FIELDS", "T_TAG", "T_INFO", "T_KEYS", "T_KEY", "T_WITH", 
	"T_VALUES", "T_VALUE", "T_FROM", "T_WHERE", "T_LIMIT", "T_QUERIES", "T_QUERY", 
	"T_EXPLAIN", "T_WITH_VALUE", "T_SELECT", "T_AS", "T_AND", "T_OR", "T_FILL", 
	"T_NULL", "T_PREVIOUS", "T_ORDER", "T_ASC", "T_DESC", "T_LIKE", "T_NOT", 
	"T_BETWEEN", "T_IS", "T_GROUP", "T_HAVING", "T_BY", "T_FOR", "T_STATS", 
	"T_TIME", "T_NOW", "T_IN", "T_LOG", "T_PROFILE", "T_SUM", "T_MIN", "T_MAX", 
	"T_COUNT", "T_AVG", "T_STDDEV", "T_QUANTILE", "T_SECOND", "T_MINUTE", "T_HOUR", 
	"T_DAY", "T_WEEK", "T_MONTH", "T_YEAR", "T_DOT", "T_COLON", "T_EQUAL", 
	"T_NOTEQUAL", "T_NOTEQUAL2", "T_GREATER", "T_GREATEREQUAL", "T_LESS", "T_LESSEQUAL", 
	"T_REGEXP", "T_NEQREGEXP", "T_COMMA", "T_OPEN_B", "T_CLOSE_B", "T_OPEN_SB", 
	"T_CLOSE_SB", "T_OPEN_P", "T_CLOSE_P", "T_ADD", "T_SUB", "T_DIV", "T_MUL", 
	"T_MOD", "L_ID", "L_INT", "L_DEC",
}

var ruleNames = []string{
	"statement", "statementList", "useStmt", "showMasterStmt", "showStoragesStmt", 
	"showAliveStmt", "createStorageStmt", "showSchemasStmt", "createDatabaseStmt", 
	"showDatabaseStmt", "showNameSpacesStmt", "showMetricsStmt", "showFieldsStmt", 
	"showTagKeysStmt", "showTagValuesStmt", "prefix", "withTagKey", "namespace", 
	"queryStmt", "selectExpr", "fields", "field", "alias", "fromClause", "whereClause", 
	"conditionExpr", "tagFilterExpr", "tagValueList", "timeRangeExpr", "timeExpr", 
	"nowExpr", "nowFunc", "groupByClause", "groupByKeys", "groupByKey", "fillOption", 
	"orderByClause", "sortField", "sortFields", "havingClause", "boolExpr", 
	"boolExprLogicalOp", "boolExprAtom", "binaryExpr", "binaryOperator", "fieldExpr", 
	"durationLit", "intervalItem", "exprFunc", "funcName", "exprFuncParams", 
	"funcParam", "exprAtom", "identFilter", "json", "obj", "pair", "arr", "value", 
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
	SQLParserT__0 = 1
	SQLParserT__1 = 2
	SQLParserT__2 = 3
	SQLParserSTRING = 4
	SQLParserWS = 5
	SQLParserT_CREATE = 6
	SQLParserT_UPDATE = 7
	SQLParserT_SET = 8
	SQLParserT_DROP = 9
	SQLParserT_INTERVAL = 10
	SQLParserT_INTERVAL_NAME = 11
	SQLParserT_SHARD = 12
	SQLParserT_REPLICATION = 13
	SQLParserT_TTL = 14
	SQLParserT_META_TTL = 15
	SQLParserT_PAST_TTL = 16
	SQLParserT_FUTURE_TTL = 17
	SQLParserT_KILL = 18
	SQLParserT_ON = 19
	SQLParserT_SHOW = 20
	SQLParserT_USE = 21
	SQLParserT_MASTER = 22
	SQLParserT_STORAGES = 23
	SQLParserT_STORAGE = 24
	SQLParserT_BROKER = 25
	SQLParserT_ALIVE = 26
	SQLParserT_SCHEMAS = 27
	SQLParserT_DATASBAE = 28
	SQLParserT_DATASBAES = 29
	SQLParserT_NAMESPACE = 30
	SQLParserT_NAMESPACES = 31
	SQLParserT_NODE = 32
	SQLParserT_METRICS = 33
	SQLParserT_METRIC = 34
	SQLParserT_FIELD = 35
	SQLParserT_FIELDS = 36
	SQLParserT_TAG = 37
	SQLParserT_INFO = 38
	SQLParserT_KEYS = 39
	SQLParserT_KEY = 40
	SQLParserT_WITH = 41
	SQLParserT_VALUES = 42
	SQLParserT_VALUE = 43
	SQLParserT_FROM = 44
	SQLParserT_WHERE = 45
	SQLParserT_LIMIT = 46
	SQLParserT_QUERIES = 47
	SQLParserT_QUERY = 48
	SQLParserT_EXPLAIN = 49
	SQLParserT_WITH_VALUE = 50
	SQLParserT_SELECT = 51
	SQLParserT_AS = 52
	SQLParserT_AND = 53
	SQLParserT_OR = 54
	SQLParserT_FILL = 55
	SQLParserT_NULL = 56
	SQLParserT_PREVIOUS = 57
	SQLParserT_ORDER = 58
	SQLParserT_ASC = 59
	SQLParserT_DESC = 60
	SQLParserT_LIKE = 61
	SQLParserT_NOT = 62
	SQLParserT_BETWEEN = 63
	SQLParserT_IS = 64
	SQLParserT_GROUP = 65
	SQLParserT_HAVING = 66
	SQLParserT_BY = 67
	SQLParserT_FOR = 68
	SQLParserT_STATS = 69
	SQLParserT_TIME = 70
	SQLParserT_NOW = 71
	SQLParserT_IN = 72
	SQLParserT_LOG = 73
	SQLParserT_PROFILE = 74
	SQLParserT_SUM = 75
	SQLParserT_MIN = 76
	SQLParserT_MAX = 77
	SQLParserT_COUNT = 78
	SQLParserT_AVG = 79
	SQLParserT_STDDEV = 80
	SQLParserT_QUANTILE = 81
	SQLParserT_SECOND = 82
	SQLParserT_MINUTE = 83
	SQLParserT_HOUR = 84
	SQLParserT_DAY = 85
	SQLParserT_WEEK = 86
	SQLParserT_MONTH = 87
	SQLParserT_YEAR = 88
	SQLParserT_DOT = 89
	SQLParserT_COLON = 90
	SQLParserT_EQUAL = 91
	SQLParserT_NOTEQUAL = 92
	SQLParserT_NOTEQUAL2 = 93
	SQLParserT_GREATER = 94
	SQLParserT_GREATEREQUAL = 95
	SQLParserT_LESS = 96
	SQLParserT_LESSEQUAL = 97
	SQLParserT_REGEXP = 98
	SQLParserT_NEQREGEXP = 99
	SQLParserT_COMMA = 100
	SQLParserT_OPEN_B = 101
	SQLParserT_CLOSE_B = 102
	SQLParserT_OPEN_SB = 103
	SQLParserT_CLOSE_SB = 104
	SQLParserT_OPEN_P = 105
	SQLParserT_CLOSE_P = 106
	SQLParserT_ADD = 107
	SQLParserT_SUB = 108
	SQLParserT_DIV = 109
	SQLParserT_MUL = 110
	SQLParserT_MOD = 111
	SQLParserL_ID = 112
	SQLParserL_INT = 113
	SQLParserL_DEC = 114
)

// SQLParser rules.
const (
	SQLParserRULE_statement = 0
	SQLParserRULE_statementList = 1
	SQLParserRULE_useStmt = 2
	SQLParserRULE_showMasterStmt = 3
	SQLParserRULE_showStoragesStmt = 4
	SQLParserRULE_showAliveStmt = 5
	SQLParserRULE_createStorageStmt = 6
	SQLParserRULE_showSchemasStmt = 7
	SQLParserRULE_createDatabaseStmt = 8
	SQLParserRULE_showDatabaseStmt = 9
	SQLParserRULE_showNameSpacesStmt = 10
	SQLParserRULE_showMetricsStmt = 11
	SQLParserRULE_showFieldsStmt = 12
	SQLParserRULE_showTagKeysStmt = 13
	SQLParserRULE_showTagValuesStmt = 14
	SQLParserRULE_prefix = 15
	SQLParserRULE_withTagKey = 16
	SQLParserRULE_namespace = 17
	SQLParserRULE_queryStmt = 18
	SQLParserRULE_selectExpr = 19
	SQLParserRULE_fields = 20
	SQLParserRULE_field = 21
	SQLParserRULE_alias = 22
	SQLParserRULE_fromClause = 23
	SQLParserRULE_whereClause = 24
	SQLParserRULE_conditionExpr = 25
	SQLParserRULE_tagFilterExpr = 26
	SQLParserRULE_tagValueList = 27
	SQLParserRULE_timeRangeExpr = 28
	SQLParserRULE_timeExpr = 29
	SQLParserRULE_nowExpr = 30
	SQLParserRULE_nowFunc = 31
	SQLParserRULE_groupByClause = 32
	SQLParserRULE_groupByKeys = 33
	SQLParserRULE_groupByKey = 34
	SQLParserRULE_fillOption = 35
	SQLParserRULE_orderByClause = 36
	SQLParserRULE_sortField = 37
	SQLParserRULE_sortFields = 38
	SQLParserRULE_havingClause = 39
	SQLParserRULE_boolExpr = 40
	SQLParserRULE_boolExprLogicalOp = 41
	SQLParserRULE_boolExprAtom = 42
	SQLParserRULE_binaryExpr = 43
	SQLParserRULE_binaryOperator = 44
	SQLParserRULE_fieldExpr = 45
	SQLParserRULE_durationLit = 46
	SQLParserRULE_intervalItem = 47
	SQLParserRULE_exprFunc = 48
	SQLParserRULE_funcName = 49
	SQLParserRULE_exprFuncParams = 50
	SQLParserRULE_funcParam = 51
	SQLParserRULE_exprAtom = 52
	SQLParserRULE_identFilter = 53
	SQLParserRULE_json = 54
	SQLParserRULE_obj = 55
	SQLParserRULE_pair = 56
	SQLParserRULE_arr = 57
	SQLParserRULE_value = 58
	SQLParserRULE_intNumber = 59
	SQLParserRULE_decNumber = 60
	SQLParserRULE_limitClause = 61
	SQLParserRULE_metricName = 62
	SQLParserRULE_tagKey = 63
	SQLParserRULE_tagValue = 64
	SQLParserRULE_ident = 65
	SQLParserRULE_nonReservedWords = 66
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
		p.SetState(134)
		p.StatementList()
	}
	{
		p.SetState(135)
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

func (s *StatementListContext) ShowStoragesStmt() IShowStoragesStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowStoragesStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowStoragesStmtContext)
}

func (s *StatementListContext) ShowAliveStmt() IShowAliveStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowAliveStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowAliveStmtContext)
}

func (s *StatementListContext) CreateStorageStmt() ICreateStorageStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICreateStorageStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICreateStorageStmtContext)
}

func (s *StatementListContext) ShowSchemasStmt() IShowSchemasStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowSchemasStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowSchemasStmtContext)
}

func (s *StatementListContext) ShowDatabaseStmt() IShowDatabaseStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowDatabaseStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowDatabaseStmtContext)
}

func (s *StatementListContext) UseStmt() IUseStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUseStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUseStmtContext)
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

func (s *StatementListContext) CreateDatabaseStmt() ICreateDatabaseStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICreateDatabaseStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICreateDatabaseStmtContext)
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

	p.SetState(151)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(137)
			p.ShowMasterStmt()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(138)
			p.ShowStoragesStmt()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(139)
			p.ShowAliveStmt()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(140)
			p.CreateStorageStmt()
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(141)
			p.ShowSchemasStmt()
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(142)
			p.ShowDatabaseStmt()
		}


	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(143)
			p.UseStmt()
		}


	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(144)
			p.ShowNameSpacesStmt()
		}


	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(145)
			p.ShowMetricsStmt()
		}


	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(146)
			p.ShowFieldsStmt()
		}


	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(147)
			p.ShowTagKeysStmt()
		}


	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(148)
			p.ShowTagValuesStmt()
		}


	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(149)
			p.QueryStmt()
		}


	case 14:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(150)
			p.CreateDatabaseStmt()
		}

	}


	return localctx
}


// IUseStmtContext is an interface to support dynamic dispatch.
type IUseStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUseStmtContext differentiates from other interfaces.
	IsUseStmtContext()
}

type UseStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUseStmtContext() *UseStmtContext {
	var p = new(UseStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_useStmt
	return p
}

func (*UseStmtContext) IsUseStmtContext() {}

func NewUseStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UseStmtContext {
	var p = new(UseStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_useStmt

	return p
}

func (s *UseStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *UseStmtContext) T_USE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_USE, 0)
}

func (s *UseStmtContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *UseStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UseStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *UseStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterUseStmt(s)
	}
}

func (s *UseStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitUseStmt(s)
	}
}




func (p *SQLParser) UseStmt() (localctx IUseStmtContext) {
	localctx = NewUseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SQLParserRULE_useStmt)

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
		p.SetState(153)
		p.Match(SQLParserT_USE)
	}
	{
		p.SetState(154)
		p.Ident()
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
	p.EnterRule(localctx, 6, SQLParserRULE_showMasterStmt)

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
		p.SetState(156)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(157)
		p.Match(SQLParserT_MASTER)
	}



	return localctx
}


// IShowStoragesStmtContext is an interface to support dynamic dispatch.
type IShowStoragesStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowStoragesStmtContext differentiates from other interfaces.
	IsShowStoragesStmtContext()
}

type ShowStoragesStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowStoragesStmtContext() *ShowStoragesStmtContext {
	var p = new(ShowStoragesStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showStoragesStmt
	return p
}

func (*ShowStoragesStmtContext) IsShowStoragesStmtContext() {}

func NewShowStoragesStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowStoragesStmtContext {
	var p = new(ShowStoragesStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showStoragesStmt

	return p
}

func (s *ShowStoragesStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowStoragesStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowStoragesStmtContext) T_STORAGES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGES, 0)
}

func (s *ShowStoragesStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStoragesStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowStoragesStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowStoragesStmt(s)
	}
}

func (s *ShowStoragesStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowStoragesStmt(s)
	}
}




func (p *SQLParser) ShowStoragesStmt() (localctx IShowStoragesStmtContext) {
	localctx = NewShowStoragesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SQLParserRULE_showStoragesStmt)

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
		p.SetState(159)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(160)
		p.Match(SQLParserT_STORAGES)
	}



	return localctx
}


// IShowAliveStmtContext is an interface to support dynamic dispatch.
type IShowAliveStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowAliveStmtContext differentiates from other interfaces.
	IsShowAliveStmtContext()
}

type ShowAliveStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowAliveStmtContext() *ShowAliveStmtContext {
	var p = new(ShowAliveStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showAliveStmt
	return p
}

func (*ShowAliveStmtContext) IsShowAliveStmtContext() {}

func NewShowAliveStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowAliveStmtContext {
	var p = new(ShowAliveStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showAliveStmt

	return p
}

func (s *ShowAliveStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowAliveStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowAliveStmtContext) T_ALIVE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ALIVE, 0)
}

func (s *ShowAliveStmtContext) T_BROKER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BROKER, 0)
}

func (s *ShowAliveStmtContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *ShowAliveStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowAliveStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowAliveStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowAliveStmt(s)
	}
}

func (s *ShowAliveStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowAliveStmt(s)
	}
}




func (p *SQLParser) ShowAliveStmt() (localctx IShowAliveStmtContext) {
	localctx = NewShowAliveStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SQLParserRULE_showAliveStmt)
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
		p.SetState(162)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(163)
		p.Match(SQLParserT_ALIVE)
	}
	{
		p.SetState(164)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserT_STORAGE || _la == SQLParserT_BROKER) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}



	return localctx
}


// ICreateStorageStmtContext is an interface to support dynamic dispatch.
type ICreateStorageStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCreateStorageStmtContext differentiates from other interfaces.
	IsCreateStorageStmtContext()
}

type CreateStorageStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateStorageStmtContext() *CreateStorageStmtContext {
	var p = new(CreateStorageStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_createStorageStmt
	return p
}

func (*CreateStorageStmtContext) IsCreateStorageStmtContext() {}

func NewCreateStorageStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateStorageStmtContext {
	var p = new(CreateStorageStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_createStorageStmt

	return p
}

func (s *CreateStorageStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateStorageStmtContext) T_CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CREATE, 0)
}

func (s *CreateStorageStmtContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *CreateStorageStmtContext) Json() IJsonContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IJsonContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IJsonContext)
}

func (s *CreateStorageStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateStorageStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *CreateStorageStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterCreateStorageStmt(s)
	}
}

func (s *CreateStorageStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitCreateStorageStmt(s)
	}
}




func (p *SQLParser) CreateStorageStmt() (localctx ICreateStorageStmtContext) {
	localctx = NewCreateStorageStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SQLParserRULE_createStorageStmt)

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
		p.SetState(166)
		p.Match(SQLParserT_CREATE)
	}
	{
		p.SetState(167)
		p.Match(SQLParserT_STORAGE)
	}
	{
		p.SetState(168)
		p.Json()
	}



	return localctx
}


// IShowSchemasStmtContext is an interface to support dynamic dispatch.
type IShowSchemasStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowSchemasStmtContext differentiates from other interfaces.
	IsShowSchemasStmtContext()
}

type ShowSchemasStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowSchemasStmtContext() *ShowSchemasStmtContext {
	var p = new(ShowSchemasStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showSchemasStmt
	return p
}

func (*ShowSchemasStmtContext) IsShowSchemasStmtContext() {}

func NewShowSchemasStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowSchemasStmtContext {
	var p = new(ShowSchemasStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showSchemasStmt

	return p
}

func (s *ShowSchemasStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowSchemasStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowSchemasStmtContext) T_SCHEMAS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SCHEMAS, 0)
}

func (s *ShowSchemasStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowSchemasStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowSchemasStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowSchemasStmt(s)
	}
}

func (s *ShowSchemasStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowSchemasStmt(s)
	}
}




func (p *SQLParser) ShowSchemasStmt() (localctx IShowSchemasStmtContext) {
	localctx = NewShowSchemasStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SQLParserRULE_showSchemasStmt)

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
		p.SetState(170)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(171)
		p.Match(SQLParserT_SCHEMAS)
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

func (s *CreateDatabaseStmtContext) Json() IJsonContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IJsonContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IJsonContext)
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




func (p *SQLParser) CreateDatabaseStmt() (localctx ICreateDatabaseStmtContext) {
	localctx = NewCreateDatabaseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SQLParserRULE_createDatabaseStmt)

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
		p.SetState(173)
		p.Match(SQLParserT_CREATE)
	}
	{
		p.SetState(174)
		p.Match(SQLParserT_DATASBAE)
	}
	{
		p.SetState(175)
		p.Json()
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
	p.EnterRule(localctx, 18, SQLParserRULE_showDatabaseStmt)

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
		p.SetState(177)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(178)
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
	p.EnterRule(localctx, 20, SQLParserRULE_showNameSpacesStmt)
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
		p.SetState(180)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(181)
		p.Match(SQLParserT_NAMESPACES)
	}
	p.SetState(186)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(182)
			p.Match(SQLParserT_WHERE)
		}
		{
			p.SetState(183)
			p.Match(SQLParserT_NAMESPACE)
		}
		{
			p.SetState(184)
			p.Match(SQLParserT_EQUAL)
		}
		{
			p.SetState(185)
			p.Prefix()
		}

	}
	p.SetState(189)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(188)
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
	p.EnterRule(localctx, 22, SQLParserRULE_showMetricsStmt)
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
		p.SetState(191)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(192)
		p.Match(SQLParserT_METRICS)
	}
	p.SetState(195)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(193)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(194)
			p.Namespace()
		}

	}
	p.SetState(201)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(197)
			p.Match(SQLParserT_WHERE)
		}
		{
			p.SetState(198)
			p.Match(SQLParserT_METRIC)
		}
		{
			p.SetState(199)
			p.Match(SQLParserT_EQUAL)
		}
		{
			p.SetState(200)
			p.Prefix()
		}

	}
	p.SetState(204)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(203)
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
	p.EnterRule(localctx, 24, SQLParserRULE_showFieldsStmt)
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
		p.SetState(206)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(207)
		p.Match(SQLParserT_FIELDS)
	}
	p.SetState(210)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(208)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(209)
			p.Namespace()
		}

	}
	{
		p.SetState(212)
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
	p.EnterRule(localctx, 26, SQLParserRULE_showTagKeysStmt)
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
		p.SetState(214)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(215)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(216)
		p.Match(SQLParserT_KEYS)
	}
	p.SetState(219)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(217)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(218)
			p.Namespace()
		}

	}
	{
		p.SetState(221)
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
	p.EnterRule(localctx, 28, SQLParserRULE_showTagValuesStmt)
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
		p.SetState(223)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(224)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(225)
		p.Match(SQLParserT_VALUES)
	}
	p.SetState(228)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(226)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(227)
			p.Namespace()
		}

	}
	{
		p.SetState(230)
		p.FromClause()
	}
	{
		p.SetState(231)
		p.Match(SQLParserT_WITH)
	}
	{
		p.SetState(232)
		p.Match(SQLParserT_KEY)
	}
	{
		p.SetState(233)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(234)
		p.WithTagKey()
	}
	p.SetState(236)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(235)
			p.WhereClause()
		}

	}
	p.SetState(239)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(238)
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
	p.EnterRule(localctx, 30, SQLParserRULE_prefix)

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
		p.SetState(241)
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
	p.EnterRule(localctx, 32, SQLParserRULE_withTagKey)

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
	p.EnterRule(localctx, 34, SQLParserRULE_namespace)

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
		p.SetState(245)
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
	p.EnterRule(localctx, 36, SQLParserRULE_queryStmt)
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
	p.SetState(248)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_EXPLAIN {
		{
			p.SetState(247)
			p.Match(SQLParserT_EXPLAIN)
		}

	}
	{
		p.SetState(250)
		p.SelectExpr()
	}
	p.SetState(253)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(251)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(252)
			p.Namespace()
		}

	}
	{
		p.SetState(255)
		p.FromClause()
	}
	p.SetState(257)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(256)
			p.WhereClause()
		}

	}
	p.SetState(260)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_GROUP {
		{
			p.SetState(259)
			p.GroupByClause()
		}

	}
	p.SetState(263)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ORDER {
		{
			p.SetState(262)
			p.OrderByClause()
		}

	}
	p.SetState(266)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(265)
			p.LimitClause()
		}

	}
	p.SetState(269)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WITH_VALUE {
		{
			p.SetState(268)
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
	p.EnterRule(localctx, 38, SQLParserRULE_selectExpr)

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
		p.SetState(271)
		p.Match(SQLParserT_SELECT)
	}
	{
		p.SetState(272)
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
	p.EnterRule(localctx, 40, SQLParserRULE_fields)
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
		p.SetState(274)
		p.Field()
	}
	p.SetState(279)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(275)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(276)
			p.Field()
		}


		p.SetState(281)
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
	p.EnterRule(localctx, 42, SQLParserRULE_field)
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
		p.SetState(282)
		p.fieldExpr(0)
	}
	p.SetState(284)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_AS {
		{
			p.SetState(283)
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
	p.EnterRule(localctx, 44, SQLParserRULE_alias)

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
		p.SetState(286)
		p.Match(SQLParserT_AS)
	}
	{
		p.SetState(287)
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
	p.EnterRule(localctx, 46, SQLParserRULE_fromClause)

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
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(290)
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
	p.EnterRule(localctx, 48, SQLParserRULE_whereClause)

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
		p.SetState(292)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(293)
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
	p.EnterRule(localctx, 50, SQLParserRULE_conditionExpr)
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

	p.SetState(305)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(295)
			p.tagFilterExpr(0)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(296)
			p.tagFilterExpr(0)
		}
		{
			p.SetState(297)
			p.Match(SQLParserT_AND)
		}
		{
			p.SetState(298)
			p.TimeRangeExpr()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(300)
			p.TimeRangeExpr()
		}
		p.SetState(303)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == SQLParserT_AND {
			{
				p.SetState(301)
				p.Match(SQLParserT_AND)
			}
			{
				p.SetState(302)
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
	_startState := 52
	p.EnterRecursionRule(localctx, 52, SQLParserRULE_tagFilterExpr, _p)
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
	p.SetState(335)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(308)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(309)
			p.tagFilterExpr(0)
		}
		{
			p.SetState(310)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(312)
			p.TagKey()
		}
		p.SetState(321)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_EQUAL:
			{
				p.SetState(313)
				p.Match(SQLParserT_EQUAL)
			}


		case SQLParserT_LIKE:
			{
				p.SetState(314)
				p.Match(SQLParserT_LIKE)
			}


		case SQLParserT_NOT:
			{
				p.SetState(315)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(316)
				p.Match(SQLParserT_LIKE)
			}


		case SQLParserT_REGEXP:
			{
				p.SetState(317)
				p.Match(SQLParserT_REGEXP)
			}


		case SQLParserT_NEQREGEXP:
			{
				p.SetState(318)
				p.Match(SQLParserT_NEQREGEXP)
			}


		case SQLParserT_NOTEQUAL:
			{
				p.SetState(319)
				p.Match(SQLParserT_NOTEQUAL)
			}


		case SQLParserT_NOTEQUAL2:
			{
				p.SetState(320)
				p.Match(SQLParserT_NOTEQUAL2)
			}



		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(323)
			p.TagValue()
		}


	case 3:
		{
			p.SetState(325)
			p.TagKey()
		}
		p.SetState(329)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_IN:
			{
				p.SetState(326)
				p.Match(SQLParserT_IN)
			}


		case SQLParserT_NOT:
			{
				p.SetState(327)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(328)
				p.Match(SQLParserT_IN)
			}



		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(331)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(332)
			p.TagValueList()
		}
		{
			p.SetState(333)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(342)
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
			p.SetState(337)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(338)
				_la = p.GetTokenStream().LA(1)

				if !(_la == SQLParserT_AND || _la == SQLParserT_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(339)
				p.tagFilterExpr(2)
			}


		}
		p.SetState(344)
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
	p.EnterRule(localctx, 54, SQLParserRULE_tagValueList)
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
		p.SetState(345)
		p.TagValue()
	}
	p.SetState(350)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(346)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(347)
			p.TagValue()
		}


		p.SetState(352)
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
	p.EnterRule(localctx, 56, SQLParserRULE_timeRangeExpr)

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
		p.SetState(353)
		p.TimeExpr()
	}
	p.SetState(356)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(354)
			p.Match(SQLParserT_AND)
		}
		{
			p.SetState(355)
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
	p.EnterRule(localctx, 58, SQLParserRULE_timeExpr)

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
		p.Match(SQLParserT_TIME)
	}
	{
		p.SetState(359)
		p.BinaryOperator()
	}
	p.SetState(362)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(360)
			p.NowExpr()
		}


	case 2:
		{
			p.SetState(361)
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
	p.EnterRule(localctx, 60, SQLParserRULE_nowExpr)
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
		p.SetState(364)
		p.NowFunc()
	}
	p.SetState(366)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if ((((_la - 107)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 107))) & ((1 << (SQLParserT_ADD - 107)) | (1 << (SQLParserT_SUB - 107)) | (1 << (SQLParserL_INT - 107)))) != 0) {
		{
			p.SetState(365)
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
	p.EnterRule(localctx, 62, SQLParserRULE_nowFunc)
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
		p.Match(SQLParserT_NOW)
	}
	{
		p.SetState(369)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(371)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_USE) | (1 << SQLParserT_MASTER) | (1 << SQLParserT_STORAGES) | (1 << SQLParserT_STORAGE) | (1 << SQLParserT_BROKER) | (1 << SQLParserT_ALIVE) | (1 << SQLParserT_SCHEMAS) | (1 << SQLParserT_DATASBAE) | (1 << SQLParserT_DATASBAES) | (1 << SQLParserT_NAMESPACE) | (1 << SQLParserT_NAMESPACES))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_NODE - 32)) | (1 << (SQLParserT_METRICS - 32)) | (1 << (SQLParserT_METRIC - 32)) | (1 << (SQLParserT_FIELD - 32)) | (1 << (SQLParserT_FIELDS - 32)) | (1 << (SQLParserT_TAG - 32)) | (1 << (SQLParserT_INFO - 32)) | (1 << (SQLParserT_KEYS - 32)) | (1 << (SQLParserT_KEY - 32)) | (1 << (SQLParserT_WITH - 32)) | (1 << (SQLParserT_VALUES - 32)) | (1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)) | (1 << (SQLParserT_ASC - 32)) | (1 << (SQLParserT_DESC - 32)) | (1 << (SQLParserT_LIKE - 32)) | (1 << (SQLParserT_NOT - 32)) | (1 << (SQLParserT_BETWEEN - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_IS - 64)) | (1 << (SQLParserT_GROUP - 64)) | (1 << (SQLParserT_HAVING - 64)) | (1 << (SQLParserT_BY - 64)) | (1 << (SQLParserT_FOR - 64)) | (1 << (SQLParserT_STATS - 64)) | (1 << (SQLParserT_TIME - 64)) | (1 << (SQLParserT_NOW - 64)) | (1 << (SQLParserT_IN - 64)) | (1 << (SQLParserT_LOG - 64)) | (1 << (SQLParserT_PROFILE - 64)) | (1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)))) != 0) || ((((_la - 105)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 105))) & ((1 << (SQLParserT_OPEN_P - 105)) | (1 << (SQLParserT_ADD - 105)) | (1 << (SQLParserT_SUB - 105)) | (1 << (SQLParserL_ID - 105)) | (1 << (SQLParserL_INT - 105)) | (1 << (SQLParserL_DEC - 105)))) != 0) {
		{
			p.SetState(370)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(373)
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
	p.EnterRule(localctx, 64, SQLParserRULE_groupByClause)
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
		p.SetState(375)
		p.Match(SQLParserT_GROUP)
	}
	{
		p.SetState(376)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(377)
		p.GroupByKeys()
	}
	p.SetState(383)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_FILL {
		{
			p.SetState(378)
			p.Match(SQLParserT_FILL)
		}
		{
			p.SetState(379)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(380)
			p.FillOption()
		}
		{
			p.SetState(381)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.SetState(386)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_HAVING {
		{
			p.SetState(385)
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
	p.EnterRule(localctx, 66, SQLParserRULE_groupByKeys)
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
		p.SetState(388)
		p.GroupByKey()
	}
	p.SetState(393)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(389)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(390)
			p.GroupByKey()
		}


		p.SetState(395)
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
	p.EnterRule(localctx, 68, SQLParserRULE_groupByKey)

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

	p.SetState(402)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(396)
			p.Ident()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(397)
			p.Match(SQLParserT_TIME)
		}
		{
			p.SetState(398)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(399)
			p.DurationLit()
		}
		{
			p.SetState(400)
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
	p.EnterRule(localctx, 70, SQLParserRULE_fillOption)
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
		p.SetState(404)
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
	p.EnterRule(localctx, 72, SQLParserRULE_orderByClause)

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
		p.SetState(406)
		p.Match(SQLParserT_ORDER)
	}
	{
		p.SetState(407)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(408)
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
	p.EnterRule(localctx, 74, SQLParserRULE_sortField)
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
		p.SetState(410)
		p.fieldExpr(0)
	}
	p.SetState(414)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_ASC || _la == SQLParserT_DESC {
		{
			p.SetState(411)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_ASC || _la == SQLParserT_DESC) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}


		p.SetState(416)
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
	p.EnterRule(localctx, 76, SQLParserRULE_sortFields)
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
		p.SetState(417)
		p.SortField()
	}
	p.SetState(422)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(418)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(419)
			p.SortField()
		}


		p.SetState(424)
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
	p.EnterRule(localctx, 78, SQLParserRULE_havingClause)

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
		p.SetState(425)
		p.Match(SQLParserT_HAVING)
	}
	{
		p.SetState(426)
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
	_startState := 80
	p.EnterRecursionRule(localctx, 80, SQLParserRULE_boolExpr, _p)

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
	p.SetState(434)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 37, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(429)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(430)
			p.boolExpr(0)
		}
		{
			p.SetState(431)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(433)
			p.BoolExprAtom()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(442)
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
			p.SetState(436)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(437)
				p.BoolExprLogicalOp()
			}
			{
				p.SetState(438)
				p.boolExpr(3)
			}


		}
		p.SetState(444)
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
	p.EnterRule(localctx, 82, SQLParserRULE_boolExprLogicalOp)
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
		p.SetState(445)
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
	p.EnterRule(localctx, 84, SQLParserRULE_boolExprAtom)

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
		p.SetState(447)
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
	p.EnterRule(localctx, 86, SQLParserRULE_binaryExpr)

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
		p.SetState(449)
		p.fieldExpr(0)
	}
	{
		p.SetState(450)
		p.BinaryOperator()
	}
	{
		p.SetState(451)
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
	p.EnterRule(localctx, 88, SQLParserRULE_binaryOperator)
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

	p.SetState(461)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_EQUAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(453)
			p.Match(SQLParserT_EQUAL)
		}


	case SQLParserT_NOTEQUAL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(454)
			p.Match(SQLParserT_NOTEQUAL)
		}


	case SQLParserT_NOTEQUAL2:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(455)
			p.Match(SQLParserT_NOTEQUAL2)
		}


	case SQLParserT_LESS:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(456)
			p.Match(SQLParserT_LESS)
		}


	case SQLParserT_LESSEQUAL:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(457)
			p.Match(SQLParserT_LESSEQUAL)
		}


	case SQLParserT_GREATER:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(458)
			p.Match(SQLParserT_GREATER)
		}


	case SQLParserT_GREATEREQUAL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(459)
			p.Match(SQLParserT_GREATEREQUAL)
		}


	case SQLParserT_LIKE, SQLParserT_REGEXP:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(460)
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
	_startState := 90
	p.EnterRecursionRule(localctx, 90, SQLParserRULE_fieldExpr, _p)

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
	p.SetState(471)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(464)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(465)
			p.fieldExpr(0)
		}
		{
			p.SetState(466)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(468)
			p.ExprFunc()
		}


	case 3:
		{
			p.SetState(469)
			p.ExprAtom()
		}


	case 4:
		{
			p.SetState(470)
			p.DurationLit()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(487)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(485)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext()) {
			case 1:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(473)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(474)
					p.Match(SQLParserT_MUL)
				}
				{
					p.SetState(475)
					p.fieldExpr(9)
				}


			case 2:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(476)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(477)
					p.Match(SQLParserT_DIV)
				}
				{
					p.SetState(478)
					p.fieldExpr(8)
				}


			case 3:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(479)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				}
				{
					p.SetState(480)
					p.Match(SQLParserT_ADD)
				}
				{
					p.SetState(481)
					p.fieldExpr(7)
				}


			case 4:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(482)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				{
					p.SetState(483)
					p.Match(SQLParserT_SUB)
				}
				{
					p.SetState(484)
					p.fieldExpr(6)
				}

			}

		}
		p.SetState(489)
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
	p.EnterRule(localctx, 92, SQLParserRULE_durationLit)

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
		p.SetState(490)
		p.IntNumber()
	}
	{
		p.SetState(491)
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
	p.EnterRule(localctx, 94, SQLParserRULE_intervalItem)
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
		p.SetState(493)
		_la = p.GetTokenStream().LA(1)

		if !(((((_la - 82)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 82))) & ((1 << (SQLParserT_SECOND - 82)) | (1 << (SQLParserT_MINUTE - 82)) | (1 << (SQLParserT_HOUR - 82)) | (1 << (SQLParserT_DAY - 82)) | (1 << (SQLParserT_WEEK - 82)) | (1 << (SQLParserT_MONTH - 82)) | (1 << (SQLParserT_YEAR - 82)))) != 0)) {
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
	p.EnterRule(localctx, 96, SQLParserRULE_exprFunc)
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
		p.SetState(495)
		p.FuncName()
	}
	{
		p.SetState(496)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(498)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_USE) | (1 << SQLParserT_MASTER) | (1 << SQLParserT_STORAGES) | (1 << SQLParserT_STORAGE) | (1 << SQLParserT_BROKER) | (1 << SQLParserT_ALIVE) | (1 << SQLParserT_SCHEMAS) | (1 << SQLParserT_DATASBAE) | (1 << SQLParserT_DATASBAES) | (1 << SQLParserT_NAMESPACE) | (1 << SQLParserT_NAMESPACES))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_NODE - 32)) | (1 << (SQLParserT_METRICS - 32)) | (1 << (SQLParserT_METRIC - 32)) | (1 << (SQLParserT_FIELD - 32)) | (1 << (SQLParserT_FIELDS - 32)) | (1 << (SQLParserT_TAG - 32)) | (1 << (SQLParserT_INFO - 32)) | (1 << (SQLParserT_KEYS - 32)) | (1 << (SQLParserT_KEY - 32)) | (1 << (SQLParserT_WITH - 32)) | (1 << (SQLParserT_VALUES - 32)) | (1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)) | (1 << (SQLParserT_ASC - 32)) | (1 << (SQLParserT_DESC - 32)) | (1 << (SQLParserT_LIKE - 32)) | (1 << (SQLParserT_NOT - 32)) | (1 << (SQLParserT_BETWEEN - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_IS - 64)) | (1 << (SQLParserT_GROUP - 64)) | (1 << (SQLParserT_HAVING - 64)) | (1 << (SQLParserT_BY - 64)) | (1 << (SQLParserT_FOR - 64)) | (1 << (SQLParserT_STATS - 64)) | (1 << (SQLParserT_TIME - 64)) | (1 << (SQLParserT_NOW - 64)) | (1 << (SQLParserT_IN - 64)) | (1 << (SQLParserT_LOG - 64)) | (1 << (SQLParserT_PROFILE - 64)) | (1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)))) != 0) || ((((_la - 105)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 105))) & ((1 << (SQLParserT_OPEN_P - 105)) | (1 << (SQLParserT_ADD - 105)) | (1 << (SQLParserT_SUB - 105)) | (1 << (SQLParserL_ID - 105)) | (1 << (SQLParserL_INT - 105)) | (1 << (SQLParserL_DEC - 105)))) != 0) {
		{
			p.SetState(497)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(500)
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
	p.EnterRule(localctx, 98, SQLParserRULE_funcName)
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
		p.SetState(502)
		_la = p.GetTokenStream().LA(1)

		if !(((((_la - 75)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 75))) & ((1 << (SQLParserT_SUM - 75)) | (1 << (SQLParserT_MIN - 75)) | (1 << (SQLParserT_MAX - 75)) | (1 << (SQLParserT_COUNT - 75)) | (1 << (SQLParserT_AVG - 75)) | (1 << (SQLParserT_STDDEV - 75)) | (1 << (SQLParserT_QUANTILE - 75)))) != 0)) {
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
	p.EnterRule(localctx, 100, SQLParserRULE_exprFuncParams)
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
		p.SetState(504)
		p.FuncParam()
	}
	p.SetState(509)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(505)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(506)
			p.FuncParam()
		}


		p.SetState(511)
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
	p.EnterRule(localctx, 102, SQLParserRULE_funcParam)

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

	p.SetState(514)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(512)
			p.fieldExpr(0)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(513)
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
	p.EnterRule(localctx, 104, SQLParserRULE_exprAtom)

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

	p.SetState(522)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(516)
			p.Ident()
		}
		p.SetState(518)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(517)
				p.IdentFilter()
			}


		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(520)
			p.DecNumber()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(521)
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
	p.EnterRule(localctx, 106, SQLParserRULE_identFilter)

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
		p.SetState(524)
		p.Match(SQLParserT_OPEN_SB)
	}
	{
		p.SetState(525)
		p.tagFilterExpr(0)
	}
	{
		p.SetState(526)
		p.Match(SQLParserT_CLOSE_SB)
	}



	return localctx
}


// IJsonContext is an interface to support dynamic dispatch.
type IJsonContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsJsonContext differentiates from other interfaces.
	IsJsonContext()
}

type JsonContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyJsonContext() *JsonContext {
	var p = new(JsonContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_json
	return p
}

func (*JsonContext) IsJsonContext() {}

func NewJsonContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *JsonContext {
	var p = new(JsonContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_json

	return p
}

func (s *JsonContext) GetParser() antlr.Parser { return s.parser }

func (s *JsonContext) Value() IValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *JsonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JsonContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *JsonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterJson(s)
	}
}

func (s *JsonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitJson(s)
	}
}




func (p *SQLParser) Json() (localctx IJsonContext) {
	localctx = NewJsonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, SQLParserRULE_json)

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
		p.Value()
	}



	return localctx
}


// IObjContext is an interface to support dynamic dispatch.
type IObjContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsObjContext differentiates from other interfaces.
	IsObjContext()
}

type ObjContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObjContext() *ObjContext {
	var p = new(ObjContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_obj
	return p
}

func (*ObjContext) IsObjContext() {}

func NewObjContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjContext {
	var p = new(ObjContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_obj

	return p
}

func (s *ObjContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjContext) T_OPEN_B() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_B, 0)
}

func (s *ObjContext) AllPair() []IPairContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPairContext)(nil)).Elem())
	var tst = make([]IPairContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPairContext)
		}
	}

	return tst
}

func (s *ObjContext) Pair(i int) IPairContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPairContext)
}

func (s *ObjContext) T_CLOSE_B() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_B, 0)
}

func (s *ObjContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *ObjContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *ObjContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ObjContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterObj(s)
	}
}

func (s *ObjContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitObj(s)
	}
}




func (p *SQLParser) Obj() (localctx IObjContext) {
	localctx = NewObjContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, SQLParserRULE_obj)
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

	p.SetState(543)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 49, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(530)
			p.Match(SQLParserT_OPEN_B)
		}
		{
			p.SetState(531)
			p.Pair()
		}
		p.SetState(536)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for _la == SQLParserT_COMMA {
			{
				p.SetState(532)
				p.Match(SQLParserT_COMMA)
			}
			{
				p.SetState(533)
				p.Pair()
			}


			p.SetState(538)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(539)
			p.Match(SQLParserT_CLOSE_B)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(541)
			p.Match(SQLParserT_OPEN_B)
		}
		{
			p.SetState(542)
			p.Match(SQLParserT_CLOSE_B)
		}

	}


	return localctx
}


// IPairContext is an interface to support dynamic dispatch.
type IPairContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPairContext differentiates from other interfaces.
	IsPairContext()
}

type PairContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPairContext() *PairContext {
	var p = new(PairContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_pair
	return p
}

func (*PairContext) IsPairContext() {}

func NewPairContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PairContext {
	var p = new(PairContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_pair

	return p
}

func (s *PairContext) GetParser() antlr.Parser { return s.parser }

func (s *PairContext) STRING() antlr.TerminalNode {
	return s.GetToken(SQLParserSTRING, 0)
}

func (s *PairContext) T_COLON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_COLON, 0)
}

func (s *PairContext) Value() IValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *PairContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PairContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterPair(s)
	}
}

func (s *PairContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitPair(s)
	}
}




func (p *SQLParser) Pair() (localctx IPairContext) {
	localctx = NewPairContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, SQLParserRULE_pair)

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
		p.Match(SQLParserSTRING)
	}
	{
		p.SetState(546)
		p.Match(SQLParserT_COLON)
	}
	{
		p.SetState(547)
		p.Value()
	}



	return localctx
}


// IArrContext is an interface to support dynamic dispatch.
type IArrContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArrContext differentiates from other interfaces.
	IsArrContext()
}

type ArrContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrContext() *ArrContext {
	var p = new(ArrContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_arr
	return p
}

func (*ArrContext) IsArrContext() {}

func NewArrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrContext {
	var p = new(ArrContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_arr

	return p
}

func (s *ArrContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrContext) T_OPEN_SB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_SB, 0)
}

func (s *ArrContext) AllValue() []IValueContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IValueContext)(nil)).Elem())
	var tst = make([]IValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IValueContext)
		}
	}

	return tst
}

func (s *ArrContext) Value(i int) IValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *ArrContext) T_CLOSE_SB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_SB, 0)
}

func (s *ArrContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *ArrContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *ArrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ArrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterArr(s)
	}
}

func (s *ArrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitArr(s)
	}
}




func (p *SQLParser) Arr() (localctx IArrContext) {
	localctx = NewArrContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, SQLParserRULE_arr)
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

	p.SetState(562)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(549)
			p.Match(SQLParserT_OPEN_SB)
		}
		{
			p.SetState(550)
			p.Value()
		}
		p.SetState(555)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for _la == SQLParserT_COMMA {
			{
				p.SetState(551)
				p.Match(SQLParserT_COMMA)
			}
			{
				p.SetState(552)
				p.Value()
			}


			p.SetState(557)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(558)
			p.Match(SQLParserT_CLOSE_SB)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(560)
			p.Match(SQLParserT_OPEN_SB)
		}
		{
			p.SetState(561)
			p.Match(SQLParserT_CLOSE_SB)
		}

	}


	return localctx
}


// IValueContext is an interface to support dynamic dispatch.
type IValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsValueContext differentiates from other interfaces.
	IsValueContext()
}

type ValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueContext() *ValueContext {
	var p = new(ValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_value
	return p
}

func (*ValueContext) IsValueContext() {}

func NewValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueContext {
	var p = new(ValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_value

	return p
}

func (s *ValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueContext) STRING() antlr.TerminalNode {
	return s.GetToken(SQLParserSTRING, 0)
}

func (s *ValueContext) IntNumber() IIntNumberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntNumberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntNumberContext)
}

func (s *ValueContext) DecNumber() IDecNumberContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDecNumberContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDecNumberContext)
}

func (s *ValueContext) Obj() IObjContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IObjContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IObjContext)
}

func (s *ValueContext) Arr() IArrContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrContext)
}

func (s *ValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterValue(s)
	}
}

func (s *ValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitValue(s)
	}
}




func (p *SQLParser) Value() (localctx IValueContext) {
	localctx = NewValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, SQLParserRULE_value)

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

	p.SetState(572)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 52, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(564)
			p.Match(SQLParserSTRING)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(565)
			p.IntNumber()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(566)
			p.DecNumber()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(567)
			p.Obj()
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(568)
			p.Arr()
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(569)
			p.Match(SQLParserT__0)
		}


	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(570)
			p.Match(SQLParserT__1)
		}


	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(571)
			p.Match(SQLParserT__2)
		}

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
	p.EnterRule(localctx, 118, SQLParserRULE_intNumber)
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
	p.SetState(575)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(574)
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
		p.SetState(577)
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
	p.EnterRule(localctx, 120, SQLParserRULE_decNumber)
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
	p.SetState(580)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(579)
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
		p.SetState(582)
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
	p.EnterRule(localctx, 122, SQLParserRULE_limitClause)

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
		p.SetState(584)
		p.Match(SQLParserT_LIMIT)
	}
	{
		p.SetState(585)
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
	p.EnterRule(localctx, 124, SQLParserRULE_metricName)

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
		p.SetState(587)
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
	p.EnterRule(localctx, 126, SQLParserRULE_tagKey)

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
		p.SetState(589)
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
	p.EnterRule(localctx, 128, SQLParserRULE_tagValue)

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
		p.SetState(591)
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
	p.EnterRule(localctx, 130, SQLParserRULE_ident)

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
	p.SetState(595)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserL_ID:
		{
			p.SetState(593)
			p.Match(SQLParserL_ID)
		}


	case SQLParserT_CREATE, SQLParserT_UPDATE, SQLParserT_SET, SQLParserT_DROP, SQLParserT_INTERVAL, SQLParserT_INTERVAL_NAME, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_META_TTL, SQLParserT_PAST_TTL, SQLParserT_FUTURE_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_USE, SQLParserT_MASTER, SQLParserT_STORAGES, SQLParserT_STORAGE, SQLParserT_BROKER, SQLParserT_ALIVE, SQLParserT_SCHEMAS, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NAMESPACE, SQLParserT_NAMESPACES, SQLParserT_NODE, SQLParserT_METRICS, SQLParserT_METRIC, SQLParserT_FIELD, SQLParserT_FIELDS, SQLParserT_TAG, SQLParserT_INFO, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_VALUE, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_EXPLAIN, SQLParserT_WITH_VALUE, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_HAVING, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_NOW, SQLParserT_IN, SQLParserT_LOG, SQLParserT_PROFILE, SQLParserT_SUM, SQLParserT_MIN, SQLParserT_MAX, SQLParserT_COUNT, SQLParserT_AVG, SQLParserT_STDDEV, SQLParserT_QUANTILE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
		{
			p.SetState(594)
			p.NonReservedWords()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(604)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 57, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(597)
				p.Match(SQLParserT_DOT)
			}
			p.SetState(600)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case SQLParserL_ID:
				{
					p.SetState(598)
					p.Match(SQLParserL_ID)
				}


			case SQLParserT_CREATE, SQLParserT_UPDATE, SQLParserT_SET, SQLParserT_DROP, SQLParserT_INTERVAL, SQLParserT_INTERVAL_NAME, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_META_TTL, SQLParserT_PAST_TTL, SQLParserT_FUTURE_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_USE, SQLParserT_MASTER, SQLParserT_STORAGES, SQLParserT_STORAGE, SQLParserT_BROKER, SQLParserT_ALIVE, SQLParserT_SCHEMAS, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NAMESPACE, SQLParserT_NAMESPACES, SQLParserT_NODE, SQLParserT_METRICS, SQLParserT_METRIC, SQLParserT_FIELD, SQLParserT_FIELDS, SQLParserT_TAG, SQLParserT_INFO, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_VALUE, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_EXPLAIN, SQLParserT_WITH_VALUE, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_HAVING, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_NOW, SQLParserT_IN, SQLParserT_LOG, SQLParserT_PROFILE, SQLParserT_SUM, SQLParserT_MIN, SQLParserT_MAX, SQLParserT_COUNT, SQLParserT_AVG, SQLParserT_STDDEV, SQLParserT_QUANTILE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
				{
					p.SetState(599)
					p.NonReservedWords()
				}



			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}


		}
		p.SetState(606)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 57, p.GetParserRuleContext())
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

func (s *NonReservedWordsContext) T_USE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_USE, 0)
}

func (s *NonReservedWordsContext) T_MASTER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MASTER, 0)
}

func (s *NonReservedWordsContext) T_STORAGES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGES, 0)
}

func (s *NonReservedWordsContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *NonReservedWordsContext) T_ALIVE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ALIVE, 0)
}

func (s *NonReservedWordsContext) T_BROKER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BROKER, 0)
}

func (s *NonReservedWordsContext) T_SCHEMAS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SCHEMAS, 0)
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
	p.EnterRule(localctx, 132, SQLParserRULE_nonReservedWords)
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
		p.SetState(607)
		_la = p.GetTokenStream().LA(1)

		if !((((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_USE) | (1 << SQLParserT_MASTER) | (1 << SQLParserT_STORAGES) | (1 << SQLParserT_STORAGE) | (1 << SQLParserT_BROKER) | (1 << SQLParserT_ALIVE) | (1 << SQLParserT_SCHEMAS) | (1 << SQLParserT_DATASBAE) | (1 << SQLParserT_DATASBAES) | (1 << SQLParserT_NAMESPACE) | (1 << SQLParserT_NAMESPACES))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_NODE - 32)) | (1 << (SQLParserT_METRICS - 32)) | (1 << (SQLParserT_METRIC - 32)) | (1 << (SQLParserT_FIELD - 32)) | (1 << (SQLParserT_FIELDS - 32)) | (1 << (SQLParserT_TAG - 32)) | (1 << (SQLParserT_INFO - 32)) | (1 << (SQLParserT_KEYS - 32)) | (1 << (SQLParserT_KEY - 32)) | (1 << (SQLParserT_WITH - 32)) | (1 << (SQLParserT_VALUES - 32)) | (1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)) | (1 << (SQLParserT_ASC - 32)) | (1 << (SQLParserT_DESC - 32)) | (1 << (SQLParserT_LIKE - 32)) | (1 << (SQLParserT_NOT - 32)) | (1 << (SQLParserT_BETWEEN - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_IS - 64)) | (1 << (SQLParserT_GROUP - 64)) | (1 << (SQLParserT_HAVING - 64)) | (1 << (SQLParserT_BY - 64)) | (1 << (SQLParserT_FOR - 64)) | (1 << (SQLParserT_STATS - 64)) | (1 << (SQLParserT_TIME - 64)) | (1 << (SQLParserT_NOW - 64)) | (1 << (SQLParserT_IN - 64)) | (1 << (SQLParserT_LOG - 64)) | (1 << (SQLParserT_PROFILE - 64)) | (1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)))) != 0)) {
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
	case 26:
			var t *TagFilterExprContext = nil
			if localctx != nil { t = localctx.(*TagFilterExprContext) }
			return p.TagFilterExpr_Sempred(t, predIndex)

	case 40:
			var t *BoolExprContext = nil
			if localctx != nil { t = localctx.(*BoolExprContext) }
			return p.BoolExpr_Sempred(t, predIndex)

	case 45:
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

