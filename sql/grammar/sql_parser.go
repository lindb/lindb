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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 123, 726, 
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
	81, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 
	3, 5, 3, 187, 10, 3, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 
	6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 
	8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10, 
	3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 5, 10, 226, 10, 10, 3, 10, 3, 10, 3, 
	10, 5, 10, 231, 10, 10, 3, 11, 3, 11, 3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 
	3, 12, 3, 12, 5, 12, 242, 10, 12, 3, 12, 3, 12, 3, 12, 5, 12, 247, 10, 
	12, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 14, 3, 14, 3, 14, 3, 14, 
	3, 14, 3, 14, 5, 14, 261, 10, 14, 3, 14, 3, 14, 3, 14, 5, 14, 266, 10, 
	14, 3, 15, 3, 15, 3, 15, 3, 15, 3, 16, 3, 16, 3, 16, 3, 17, 3, 17, 3, 17, 
	3, 17, 3, 18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 3, 19, 3, 19, 3, 19, 5, 
	19, 288, 10, 19, 3, 19, 5, 19, 291, 10, 19, 3, 20, 3, 20, 3, 20, 3, 20, 
	5, 20, 297, 10, 20, 3, 20, 3, 20, 3, 20, 3, 20, 5, 20, 303, 10, 20, 3, 
	20, 5, 20, 306, 10, 20, 3, 21, 3, 21, 3, 21, 3, 21, 3, 22, 3, 22, 3, 22, 
	3, 22, 3, 22, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 
	23, 5, 23, 326, 10, 23, 3, 23, 5, 23, 329, 10, 23, 3, 24, 3, 24, 3, 25, 
	3, 25, 3, 26, 3, 26, 3, 27, 3, 27, 3, 28, 5, 28, 340, 10, 28, 3, 28, 3, 
	28, 3, 28, 5, 28, 345, 10, 28, 3, 28, 5, 28, 348, 10, 28, 3, 28, 5, 28, 
	351, 10, 28, 3, 28, 5, 28, 354, 10, 28, 3, 28, 5, 28, 357, 10, 28, 3, 29, 
	3, 29, 3, 29, 3, 30, 3, 30, 3, 30, 7, 30, 365, 10, 30, 12, 30, 14, 30, 
	368, 11, 30, 3, 31, 3, 31, 5, 31, 372, 10, 31, 3, 32, 3, 32, 3, 32, 3, 
	33, 3, 33, 3, 33, 3, 33, 3, 34, 3, 34, 3, 34, 3, 34, 3, 35, 3, 35, 3, 35, 
	3, 35, 3, 36, 3, 36, 3, 36, 3, 36, 5, 36, 393, 10, 36, 3, 37, 3, 37, 3, 
	37, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 5, 38, 406, 
	10, 38, 5, 38, 408, 10, 38, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 
	39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 5, 39, 424, 10, 39, 
	3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 5, 39, 432, 10, 39, 3, 39, 3, 
	39, 3, 39, 3, 39, 5, 39, 438, 10, 39, 3, 39, 3, 39, 3, 39, 7, 39, 443, 
	10, 39, 12, 39, 14, 39, 446, 11, 39, 3, 40, 3, 40, 3, 40, 7, 40, 451, 10, 
	40, 12, 40, 14, 40, 454, 11, 40, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 
	41, 3, 42, 3, 42, 3, 42, 7, 42, 465, 10, 42, 12, 42, 14, 42, 468, 11, 42, 
	3, 43, 3, 43, 3, 43, 5, 43, 473, 10, 43, 3, 44, 3, 44, 3, 44, 3, 44, 5, 
	44, 479, 10, 44, 3, 45, 3, 45, 5, 45, 483, 10, 45, 3, 46, 3, 46, 3, 46, 
	5, 46, 488, 10, 46, 3, 46, 3, 46, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 
	47, 3, 47, 3, 47, 5, 47, 500, 10, 47, 3, 47, 5, 47, 503, 10, 47, 3, 48, 
	3, 48, 3, 48, 7, 48, 508, 10, 48, 12, 48, 14, 48, 511, 11, 48, 3, 49, 3, 
	49, 3, 49, 3, 49, 3, 49, 3, 49, 5, 49, 519, 10, 49, 3, 50, 3, 50, 3, 51, 
	3, 51, 3, 51, 3, 51, 3, 52, 3, 52, 7, 52, 529, 10, 52, 12, 52, 14, 52, 
	532, 11, 52, 3, 53, 3, 53, 3, 53, 7, 53, 537, 10, 53, 12, 53, 14, 53, 540, 
	11, 53, 3, 54, 3, 54, 3, 54, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 3, 55, 
	5, 55, 551, 10, 55, 3, 55, 3, 55, 3, 55, 3, 55, 7, 55, 557, 10, 55, 12, 
	55, 14, 55, 560, 11, 55, 3, 56, 3, 56, 3, 57, 3, 57, 3, 58, 3, 58, 3, 58, 
	3, 58, 3, 59, 3, 59, 3, 59, 3, 59, 3, 59, 3, 59, 3, 59, 3, 59, 5, 59, 578, 
	10, 59, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 5, 60, 
	588, 10, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 60, 3, 
	60, 3, 60, 3, 60, 3, 60, 7, 60, 602, 10, 60, 12, 60, 14, 60, 605, 11, 60, 
	3, 61, 3, 61, 3, 61, 3, 62, 3, 62, 3, 63, 3, 63, 3, 63, 5, 63, 615, 10, 
	63, 3, 63, 3, 63, 3, 64, 3, 64, 3, 65, 3, 65, 3, 65, 7, 65, 624, 10, 65, 
	12, 65, 14, 65, 627, 11, 65, 3, 66, 3, 66, 5, 66, 631, 10, 66, 3, 67, 3, 
	67, 5, 67, 635, 10, 67, 3, 67, 3, 67, 5, 67, 639, 10, 67, 3, 68, 3, 68, 
	3, 68, 3, 68, 3, 69, 3, 69, 3, 70, 3, 70, 3, 70, 3, 70, 7, 70, 651, 10, 
	70, 12, 70, 14, 70, 654, 11, 70, 3, 70, 3, 70, 3, 70, 3, 70, 5, 70, 660, 
	10, 70, 3, 71, 3, 71, 3, 71, 3, 71, 3, 72, 3, 72, 3, 72, 3, 72, 7, 72, 
	670, 10, 72, 12, 72, 14, 72, 673, 11, 72, 3, 72, 3, 72, 3, 72, 3, 72, 5, 
	72, 679, 10, 72, 3, 73, 3, 73, 3, 73, 3, 73, 3, 73, 3, 73, 3, 73, 3, 73, 
	5, 73, 689, 10, 73, 3, 74, 5, 74, 692, 10, 74, 3, 74, 3, 74, 3, 75, 5, 
	75, 697, 10, 75, 3, 75, 3, 75, 3, 76, 3, 76, 3, 76, 3, 77, 3, 77, 3, 78, 
	3, 78, 3, 79, 3, 79, 3, 80, 3, 80, 5, 80, 712, 10, 80, 3, 80, 3, 80, 3, 
	80, 5, 80, 717, 10, 80, 7, 80, 719, 10, 80, 12, 80, 14, 80, 722, 11, 80, 
	3, 81, 3, 81, 3, 81, 2, 5, 76, 108, 118, 82, 2, 4, 6, 8, 10, 12, 14, 16, 
	18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 
	54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 
	90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 110, 112, 114, 116, 118, 120, 
	122, 124, 126, 128, 130, 132, 134, 136, 138, 140, 142, 144, 146, 148, 150, 
	152, 154, 156, 158, 160, 2, 12, 3, 2, 31, 32, 3, 2, 24, 25, 3, 2, 60, 61, 
	4, 2, 63, 64, 122, 123, 3, 2, 66, 67, 4, 2, 68, 68, 106, 106, 3, 2, 90, 
	96, 3, 2, 82, 89, 3, 2, 115, 116, 3, 2, 8, 96, 2, 750, 2, 162, 3, 2, 2, 
	2, 4, 186, 3, 2, 2, 2, 6, 188, 3, 2, 2, 2, 8, 191, 3, 2, 2, 2, 10, 194, 
	3, 2, 2, 2, 12, 197, 3, 2, 2, 2, 14, 201, 3, 2, 2, 2, 16, 209, 3, 2, 2, 
	2, 18, 217, 3, 2, 2, 2, 20, 232, 3, 2, 2, 2, 22, 236, 3, 2, 2, 2, 24, 248, 
	3, 2, 2, 2, 26, 254, 3, 2, 2, 2, 28, 267, 3, 2, 2, 2, 30, 271, 3, 2, 2, 
	2, 32, 274, 3, 2, 2, 2, 34, 278, 3, 2, 2, 2, 36, 281, 3, 2, 2, 2, 38, 292, 
	3, 2, 2, 2, 40, 307, 3, 2, 2, 2, 42, 311, 3, 2, 2, 2, 44, 316, 3, 2, 2, 
	2, 46, 330, 3, 2, 2, 2, 48, 332, 3, 2, 2, 2, 50, 334, 3, 2, 2, 2, 52, 336, 
	3, 2, 2, 2, 54, 339, 3, 2, 2, 2, 56, 358, 3, 2, 2, 2, 58, 361, 3, 2, 2, 
	2, 60, 369, 3, 2, 2, 2, 62, 373, 3, 2, 2, 2, 64, 376, 3, 2, 2, 2, 66, 380, 
	3, 2, 2, 2, 68, 384, 3, 2, 2, 2, 70, 388, 3, 2, 2, 2, 72, 394, 3, 2, 2, 
	2, 74, 407, 3, 2, 2, 2, 76, 437, 3, 2, 2, 2, 78, 447, 3, 2, 2, 2, 80, 455, 
	3, 2, 2, 2, 82, 461, 3, 2, 2, 2, 84, 469, 3, 2, 2, 2, 86, 474, 3, 2, 2, 
	2, 88, 480, 3, 2, 2, 2, 90, 484, 3, 2, 2, 2, 92, 491, 3, 2, 2, 2, 94, 504, 
	3, 2, 2, 2, 96, 518, 3, 2, 2, 2, 98, 520, 3, 2, 2, 2, 100, 522, 3, 2, 2, 
	2, 102, 526, 3, 2, 2, 2, 104, 533, 3, 2, 2, 2, 106, 541, 3, 2, 2, 2, 108, 
	550, 3, 2, 2, 2, 110, 561, 3, 2, 2, 2, 112, 563, 3, 2, 2, 2, 114, 565, 
	3, 2, 2, 2, 116, 577, 3, 2, 2, 2, 118, 587, 3, 2, 2, 2, 120, 606, 3, 2, 
	2, 2, 122, 609, 3, 2, 2, 2, 124, 611, 3, 2, 2, 2, 126, 618, 3, 2, 2, 2, 
	128, 620, 3, 2, 2, 2, 130, 630, 3, 2, 2, 2, 132, 638, 3, 2, 2, 2, 134, 
	640, 3, 2, 2, 2, 136, 644, 3, 2, 2, 2, 138, 659, 3, 2, 2, 2, 140, 661, 
	3, 2, 2, 2, 142, 678, 3, 2, 2, 2, 144, 688, 3, 2, 2, 2, 146, 691, 3, 2, 
	2, 2, 148, 696, 3, 2, 2, 2, 150, 700, 3, 2, 2, 2, 152, 703, 3, 2, 2, 2, 
	154, 705, 3, 2, 2, 2, 156, 707, 3, 2, 2, 2, 158, 711, 3, 2, 2, 2, 160, 
	723, 3, 2, 2, 2, 162, 163, 5, 4, 3, 2, 163, 164, 7, 2, 2, 3, 164, 3, 3, 
	2, 2, 2, 165, 187, 5, 8, 5, 2, 166, 187, 5, 12, 7, 2, 167, 187, 5, 14, 
	8, 2, 168, 187, 5, 16, 9, 2, 169, 187, 5, 18, 10, 2, 170, 187, 5, 10, 6, 
	2, 171, 187, 5, 20, 11, 2, 172, 187, 5, 24, 13, 2, 173, 187, 5, 26, 14, 
	2, 174, 187, 5, 28, 15, 2, 175, 187, 5, 22, 12, 2, 176, 187, 5, 30, 16, 
	2, 177, 187, 5, 34, 18, 2, 178, 187, 5, 6, 4, 2, 179, 187, 5, 36, 19, 2, 
	180, 187, 5, 38, 20, 2, 181, 187, 5, 40, 21, 2, 182, 187, 5, 42, 22, 2, 
	183, 187, 5, 44, 23, 2, 184, 187, 5, 54, 28, 2, 185, 187, 5, 32, 17, 2, 
	186, 165, 3, 2, 2, 2, 186, 166, 3, 2, 2, 2, 186, 167, 3, 2, 2, 2, 186, 
	168, 3, 2, 2, 2, 186, 169, 3, 2, 2, 2, 186, 170, 3, 2, 2, 2, 186, 171, 
	3, 2, 2, 2, 186, 172, 3, 2, 2, 2, 186, 173, 3, 2, 2, 2, 186, 174, 3, 2, 
	2, 2, 186, 175, 3, 2, 2, 2, 186, 176, 3, 2, 2, 2, 186, 177, 3, 2, 2, 2, 
	186, 178, 3, 2, 2, 2, 186, 179, 3, 2, 2, 2, 186, 180, 3, 2, 2, 2, 186, 
	181, 3, 2, 2, 2, 186, 182, 3, 2, 2, 2, 186, 183, 3, 2, 2, 2, 186, 184, 
	3, 2, 2, 2, 186, 185, 3, 2, 2, 2, 187, 5, 3, 2, 2, 2, 188, 189, 7, 23, 
	2, 2, 189, 190, 5, 158, 80, 2, 190, 7, 3, 2, 2, 2, 191, 192, 7, 22, 2, 
	2, 192, 193, 7, 26, 2, 2, 193, 9, 3, 2, 2, 2, 194, 195, 7, 22, 2, 2, 195, 
	196, 7, 30, 2, 2, 196, 11, 3, 2, 2, 2, 197, 198, 7, 22, 2, 2, 198, 199, 
	7, 27, 2, 2, 199, 200, 7, 28, 2, 2, 200, 13, 3, 2, 2, 2, 201, 202, 7, 22, 
	2, 2, 202, 203, 7, 32, 2, 2, 203, 204, 7, 27, 2, 2, 204, 205, 7, 51, 2, 
	2, 205, 206, 5, 52, 27, 2, 206, 207, 7, 52, 2, 2, 207, 208, 5, 68, 35, 
	2, 208, 15, 3, 2, 2, 2, 209, 210, 7, 22, 2, 2, 210, 211, 7, 26, 2, 2, 211, 
	212, 7, 27, 2, 2, 212, 213, 7, 51, 2, 2, 213, 214, 5, 52, 27, 2, 214, 215, 
	7, 52, 2, 2, 215, 216, 5, 68, 35, 2, 216, 17, 3, 2, 2, 2, 217, 218, 7, 
	22, 2, 2, 218, 219, 7, 31, 2, 2, 219, 220, 7, 27, 2, 2, 220, 221, 7, 51, 
	2, 2, 221, 222, 5, 52, 27, 2, 222, 225, 7, 52, 2, 2, 223, 226, 5, 64, 33, 
	2, 224, 226, 5, 68, 35, 2, 225, 223, 3, 2, 2, 2, 225, 224, 3, 2, 2, 2, 
	226, 227, 3, 2, 2, 2, 227, 230, 7, 60, 2, 2, 228, 231, 5, 64, 33, 2, 229, 
	231, 5, 68, 35, 2, 230, 228, 3, 2, 2, 2, 230, 229, 3, 2, 2, 2, 231, 19, 
	3, 2, 2, 2, 232, 233, 7, 22, 2, 2, 233, 234, 9, 2, 2, 2, 234, 235, 7, 33, 
	2, 2, 235, 21, 3, 2, 2, 2, 236, 237, 7, 22, 2, 2, 237, 238, 7, 15, 2, 2, 
	238, 241, 7, 52, 2, 2, 239, 242, 5, 64, 33, 2, 240, 242, 5, 66, 34, 2, 
	241, 239, 3, 2, 2, 2, 241, 240, 3, 2, 2, 2, 242, 243, 3, 2, 2, 2, 243, 
	246, 7, 60, 2, 2, 244, 247, 5, 64, 33, 2, 245, 247, 5, 66, 34, 2, 246, 
	244, 3, 2, 2, 2, 246, 245, 3, 2, 2, 2, 247, 23, 3, 2, 2, 2, 248, 249, 7, 
	22, 2, 2, 249, 250, 7, 32, 2, 2, 250, 251, 7, 41, 2, 2, 251, 252, 7, 52, 
	2, 2, 252, 253, 5, 80, 41, 2, 253, 25, 3, 2, 2, 2, 254, 255, 7, 22, 2, 
	2, 255, 256, 7, 31, 2, 2, 256, 257, 7, 41, 2, 2, 257, 260, 7, 52, 2, 2, 
	258, 261, 5, 64, 33, 2, 259, 261, 5, 80, 41, 2, 260, 258, 3, 2, 2, 2, 260, 
	259, 3, 2, 2, 2, 261, 262, 3, 2, 2, 2, 262, 265, 7, 60, 2, 2, 263, 266, 
	5, 64, 33, 2, 264, 266, 5, 80, 41, 2, 265, 263, 3, 2, 2, 2, 265, 264, 3, 
	2, 2, 2, 266, 27, 3, 2, 2, 2, 267, 268, 7, 8, 2, 2, 268, 269, 7, 31, 2, 
	2, 269, 270, 5, 136, 69, 2, 270, 29, 3, 2, 2, 2, 271, 272, 7, 22, 2, 2, 
	272, 273, 7, 34, 2, 2, 273, 31, 3, 2, 2, 2, 274, 275, 7, 8, 2, 2, 275, 
	276, 7, 35, 2, 2, 276, 277, 5, 136, 69, 2, 277, 33, 3, 2, 2, 2, 278, 279, 
	7, 22, 2, 2, 279, 280, 7, 36, 2, 2, 280, 35, 3, 2, 2, 2, 281, 282, 7, 22, 
	2, 2, 282, 287, 7, 38, 2, 2, 283, 284, 7, 52, 2, 2, 284, 285, 7, 37, 2, 
	2, 285, 286, 7, 99, 2, 2, 286, 288, 5, 46, 24, 2, 287, 283, 3, 2, 2, 2, 
	287, 288, 3, 2, 2, 2, 288, 290, 3, 2, 2, 2, 289, 291, 5, 150, 76, 2, 290, 
	289, 3, 2, 2, 2, 290, 291, 3, 2, 2, 2, 291, 37, 3, 2, 2, 2, 292, 293, 7, 
	22, 2, 2, 293, 296, 7, 40, 2, 2, 294, 295, 7, 21, 2, 2, 295, 297, 5, 50, 
	26, 2, 296, 294, 3, 2, 2, 2, 296, 297, 3, 2, 2, 2, 297, 302, 3, 2, 2, 2, 
	298, 299, 7, 52, 2, 2, 299, 300, 7, 41, 2, 2, 300, 301, 7, 99, 2, 2, 301, 
	303, 5, 46, 24, 2, 302, 298, 3, 2, 2, 2, 302, 303, 3, 2, 2, 2, 303, 305, 
	3, 2, 2, 2, 304, 306, 5, 150, 76, 2, 305, 304, 3, 2, 2, 2, 305, 306, 3, 
	2, 2, 2, 306, 39, 3, 2, 2, 2, 307, 308, 7, 22, 2, 2, 308, 309, 7, 43, 2, 
	2, 309, 310, 5, 70, 36, 2, 310, 41, 3, 2, 2, 2, 311, 312, 7, 22, 2, 2, 
	312, 313, 7, 44, 2, 2, 313, 314, 7, 46, 2, 2, 314, 315, 5, 70, 36, 2, 315, 
	43, 3, 2, 2, 2, 316, 317, 7, 22, 2, 2, 317, 318, 7, 44, 2, 2, 318, 319, 
	7, 49, 2, 2, 319, 320, 5, 70, 36, 2, 320, 321, 7, 48, 2, 2, 321, 322, 7, 
	47, 2, 2, 322, 323, 7, 99, 2, 2, 323, 325, 5, 48, 25, 2, 324, 326, 5, 72, 
	37, 2, 325, 324, 3, 2, 2, 2, 325, 326, 3, 2, 2, 2, 326, 328, 3, 2, 2, 2, 
	327, 329, 5, 150, 76, 2, 328, 327, 3, 2, 2, 2, 328, 329, 3, 2, 2, 2, 329, 
	45, 3, 2, 2, 2, 330, 331, 5, 158, 80, 2, 331, 47, 3, 2, 2, 2, 332, 333, 
	5, 158, 80, 2, 333, 49, 3, 2, 2, 2, 334, 335, 5, 158, 80, 2, 335, 51, 3, 
	2, 2, 2, 336, 337, 9, 3, 2, 2, 337, 53, 3, 2, 2, 2, 338, 340, 7, 56, 2, 
	2, 339, 338, 3, 2, 2, 2, 339, 340, 3, 2, 2, 2, 340, 341, 3, 2, 2, 2, 341, 
	342, 5, 56, 29, 2, 342, 344, 5, 70, 36, 2, 343, 345, 5, 72, 37, 2, 344, 
	343, 3, 2, 2, 2, 344, 345, 3, 2, 2, 2, 345, 347, 3, 2, 2, 2, 346, 348, 
	5, 92, 47, 2, 347, 346, 3, 2, 2, 2, 347, 348, 3, 2, 2, 2, 348, 350, 3, 
	2, 2, 2, 349, 351, 5, 100, 51, 2, 350, 349, 3, 2, 2, 2, 350, 351, 3, 2, 
	2, 2, 351, 353, 3, 2, 2, 2, 352, 354, 5, 150, 76, 2, 353, 352, 3, 2, 2, 
	2, 353, 354, 3, 2, 2, 2, 354, 356, 3, 2, 2, 2, 355, 357, 7, 57, 2, 2, 356, 
	355, 3, 2, 2, 2, 356, 357, 3, 2, 2, 2, 357, 55, 3, 2, 2, 2, 358, 359, 7, 
	58, 2, 2, 359, 360, 5, 58, 30, 2, 360, 57, 3, 2, 2, 2, 361, 366, 5, 60, 
	31, 2, 362, 363, 7, 108, 2, 2, 363, 365, 5, 60, 31, 2, 364, 362, 3, 2, 
	2, 2, 365, 368, 3, 2, 2, 2, 366, 364, 3, 2, 2, 2, 366, 367, 3, 2, 2, 2, 
	367, 59, 3, 2, 2, 2, 368, 366, 3, 2, 2, 2, 369, 371, 5, 118, 60, 2, 370, 
	372, 5, 62, 32, 2, 371, 370, 3, 2, 2, 2, 371, 372, 3, 2, 2, 2, 372, 61, 
	3, 2, 2, 2, 373, 374, 7, 59, 2, 2, 374, 375, 5, 158, 80, 2, 375, 63, 3, 
	2, 2, 2, 376, 377, 7, 31, 2, 2, 377, 378, 7, 99, 2, 2, 378, 379, 5, 158, 
	80, 2, 379, 65, 3, 2, 2, 2, 380, 381, 7, 35, 2, 2, 381, 382, 7, 99, 2, 
	2, 382, 383, 5, 158, 80, 2, 383, 67, 3, 2, 2, 2, 384, 385, 7, 29, 2, 2, 
	385, 386, 7, 99, 2, 2, 386, 387, 5, 158, 80, 2, 387, 69, 3, 2, 2, 2, 388, 
	389, 7, 51, 2, 2, 389, 392, 5, 152, 77, 2, 390, 391, 7, 21, 2, 2, 391, 
	393, 5, 50, 26, 2, 392, 390, 3, 2, 2, 2, 392, 393, 3, 2, 2, 2, 393, 71, 
	3, 2, 2, 2, 394, 395, 7, 52, 2, 2, 395, 396, 5, 74, 38, 2, 396, 73, 3, 
	2, 2, 2, 397, 408, 5, 76, 39, 2, 398, 399, 5, 76, 39, 2, 399, 400, 7, 60, 
	2, 2, 400, 401, 5, 84, 43, 2, 401, 408, 3, 2, 2, 2, 402, 405, 5, 84, 43, 
	2, 403, 404, 7, 60, 2, 2, 404, 406, 5, 76, 39, 2, 405, 403, 3, 2, 2, 2, 
	405, 406, 3, 2, 2, 2, 406, 408, 3, 2, 2, 2, 407, 397, 3, 2, 2, 2, 407, 
	398, 3, 2, 2, 2, 407, 402, 3, 2, 2, 2, 408, 75, 3, 2, 2, 2, 409, 410, 8, 
	39, 1, 2, 410, 411, 7, 113, 2, 2, 411, 412, 5, 76, 39, 2, 412, 413, 7, 
	114, 2, 2, 413, 438, 3, 2, 2, 2, 414, 423, 5, 154, 78, 2, 415, 424, 7, 
	99, 2, 2, 416, 424, 7, 68, 2, 2, 417, 418, 7, 69, 2, 2, 418, 424, 7, 68, 
	2, 2, 419, 424, 7, 106, 2, 2, 420, 424, 7, 107, 2, 2, 421, 424, 7, 100, 
	2, 2, 422, 424, 7, 101, 2, 2, 423, 415, 3, 2, 2, 2, 423, 416, 3, 2, 2, 
	2, 423, 417, 3, 2, 2, 2, 423, 419, 3, 2, 2, 2, 423, 420, 3, 2, 2, 2, 423, 
	421, 3, 2, 2, 2, 423, 422, 3, 2, 2, 2, 424, 425, 3, 2, 2, 2, 425, 426, 
	5, 156, 79, 2, 426, 438, 3, 2, 2, 2, 427, 431, 5, 154, 78, 2, 428, 432, 
	7, 79, 2, 2, 429, 430, 7, 69, 2, 2, 430, 432, 7, 79, 2, 2, 431, 428, 3, 
	2, 2, 2, 431, 429, 3, 2, 2, 2, 432, 433, 3, 2, 2, 2, 433, 434, 7, 113, 
	2, 2, 434, 435, 5, 78, 40, 2, 435, 436, 7, 114, 2, 2, 436, 438, 3, 2, 2, 
	2, 437, 409, 3, 2, 2, 2, 437, 414, 3, 2, 2, 2, 437, 427, 3, 2, 2, 2, 438, 
	444, 3, 2, 2, 2, 439, 440, 12, 3, 2, 2, 440, 441, 9, 4, 2, 2, 441, 443, 
	5, 76, 39, 4, 442, 439, 3, 2, 2, 2, 443, 446, 3, 2, 2, 2, 444, 442, 3, 
	2, 2, 2, 444, 445, 3, 2, 2, 2, 445, 77, 3, 2, 2, 2, 446, 444, 3, 2, 2, 
	2, 447, 452, 5, 156, 79, 2, 448, 449, 7, 108, 2, 2, 449, 451, 5, 156, 79, 
	2, 450, 448, 3, 2, 2, 2, 451, 454, 3, 2, 2, 2, 452, 450, 3, 2, 2, 2, 452, 
	453, 3, 2, 2, 2, 453, 79, 3, 2, 2, 2, 454, 452, 3, 2, 2, 2, 455, 456, 7, 
	41, 2, 2, 456, 457, 7, 79, 2, 2, 457, 458, 7, 113, 2, 2, 458, 459, 5, 82, 
	42, 2, 459, 460, 7, 114, 2, 2, 460, 81, 3, 2, 2, 2, 461, 466, 5, 158, 80, 
	2, 462, 463, 7, 108, 2, 2, 463, 465, 5, 158, 80, 2, 464, 462, 3, 2, 2, 
	2, 465, 468, 3, 2, 2, 2, 466, 464, 3, 2, 2, 2, 466, 467, 3, 2, 2, 2, 467, 
	83, 3, 2, 2, 2, 468, 466, 3, 2, 2, 2, 469, 472, 5, 86, 44, 2, 470, 471, 
	7, 60, 2, 2, 471, 473, 5, 86, 44, 2, 472, 470, 3, 2, 2, 2, 472, 473, 3, 
	2, 2, 2, 473, 85, 3, 2, 2, 2, 474, 475, 7, 77, 2, 2, 475, 478, 5, 116, 
	59, 2, 476, 479, 5, 88, 45, 2, 477, 479, 5, 158, 80, 2, 478, 476, 3, 2, 
	2, 2, 478, 477, 3, 2, 2, 2, 479, 87, 3, 2, 2, 2, 480, 482, 5, 90, 46, 2, 
	481, 483, 5, 120, 61, 2, 482, 481, 3, 2, 2, 2, 482, 483, 3, 2, 2, 2, 483, 
	89, 3, 2, 2, 2, 484, 485, 7, 78, 2, 2, 485, 487, 7, 113, 2, 2, 486, 488, 
	5, 128, 65, 2, 487, 486, 3, 2, 2, 2, 487, 488, 3, 2, 2, 2, 488, 489, 3, 
	2, 2, 2, 489, 490, 7, 114, 2, 2, 490, 91, 3, 2, 2, 2, 491, 492, 7, 72, 
	2, 2, 492, 493, 7, 74, 2, 2, 493, 499, 5, 94, 48, 2, 494, 495, 7, 62, 2, 
	2, 495, 496, 7, 113, 2, 2, 496, 497, 5, 98, 50, 2, 497, 498, 7, 114, 2, 
	2, 498, 500, 3, 2, 2, 2, 499, 494, 3, 2, 2, 2, 499, 500, 3, 2, 2, 2, 500, 
	502, 3, 2, 2, 2, 501, 503, 5, 106, 54, 2, 502, 501, 3, 2, 2, 2, 502, 503, 
	3, 2, 2, 2, 503, 93, 3, 2, 2, 2, 504, 509, 5, 96, 49, 2, 505, 506, 7, 108, 
	2, 2, 506, 508, 5, 96, 49, 2, 507, 505, 3, 2, 2, 2, 508, 511, 3, 2, 2, 
	2, 509, 507, 3, 2, 2, 2, 509, 510, 3, 2, 2, 2, 510, 95, 3, 2, 2, 2, 511, 
	509, 3, 2, 2, 2, 512, 519, 5, 158, 80, 2, 513, 514, 7, 77, 2, 2, 514, 515, 
	7, 113, 2, 2, 515, 516, 5, 120, 61, 2, 516, 517, 7, 114, 2, 2, 517, 519, 
	3, 2, 2, 2, 518, 512, 3, 2, 2, 2, 518, 513, 3, 2, 2, 2, 519, 97, 3, 2, 
	2, 2, 520, 521, 9, 5, 2, 2, 521, 99, 3, 2, 2, 2, 522, 523, 7, 65, 2, 2, 
	523, 524, 7, 74, 2, 2, 524, 525, 5, 104, 53, 2, 525, 101, 3, 2, 2, 2, 526, 
	530, 5, 118, 60, 2, 527, 529, 9, 6, 2, 2, 528, 527, 3, 2, 2, 2, 529, 532, 
	3, 2, 2, 2, 530, 528, 3, 2, 2, 2, 530, 531, 3, 2, 2, 2, 531, 103, 3, 2, 
	2, 2, 532, 530, 3, 2, 2, 2, 533, 538, 5, 102, 52, 2, 534, 535, 7, 108, 
	2, 2, 535, 537, 5, 102, 52, 2, 536, 534, 3, 2, 2, 2, 537, 540, 3, 2, 2, 
	2, 538, 536, 3, 2, 2, 2, 538, 539, 3, 2, 2, 2, 539, 105, 3, 2, 2, 2, 540, 
	538, 3, 2, 2, 2, 541, 542, 7, 73, 2, 2, 542, 543, 5, 108, 55, 2, 543, 107, 
	3, 2, 2, 2, 544, 545, 8, 55, 1, 2, 545, 546, 7, 113, 2, 2, 546, 547, 5, 
	108, 55, 2, 547, 548, 7, 114, 2, 2, 548, 551, 3, 2, 2, 2, 549, 551, 5, 
	112, 57, 2, 550, 544, 3, 2, 2, 2, 550, 549, 3, 2, 2, 2, 551, 558, 3, 2, 
	2, 2, 552, 553, 12, 4, 2, 2, 553, 554, 5, 110, 56, 2, 554, 555, 5, 108, 
	55, 5, 555, 557, 3, 2, 2, 2, 556, 552, 3, 2, 2, 2, 557, 560, 3, 2, 2, 2, 
	558, 556, 3, 2, 2, 2, 558, 559, 3, 2, 2, 2, 559, 109, 3, 2, 2, 2, 560, 
	558, 3, 2, 2, 2, 561, 562, 9, 4, 2, 2, 562, 111, 3, 2, 2, 2, 563, 564, 
	5, 114, 58, 2, 564, 113, 3, 2, 2, 2, 565, 566, 5, 118, 60, 2, 566, 567, 
	5, 116, 59, 2, 567, 568, 5, 118, 60, 2, 568, 115, 3, 2, 2, 2, 569, 578, 
	7, 99, 2, 2, 570, 578, 7, 100, 2, 2, 571, 578, 7, 101, 2, 2, 572, 578, 
	7, 104, 2, 2, 573, 578, 7, 105, 2, 2, 574, 578, 7, 102, 2, 2, 575, 578, 
	7, 103, 2, 2, 576, 578, 9, 7, 2, 2, 577, 569, 3, 2, 2, 2, 577, 570, 3, 
	2, 2, 2, 577, 571, 3, 2, 2, 2, 577, 572, 3, 2, 2, 2, 577, 573, 3, 2, 2, 
	2, 577, 574, 3, 2, 2, 2, 577, 575, 3, 2, 2, 2, 577, 576, 3, 2, 2, 2, 578, 
	117, 3, 2, 2, 2, 579, 580, 8, 60, 1, 2, 580, 581, 7, 113, 2, 2, 581, 582, 
	5, 118, 60, 2, 582, 583, 7, 114, 2, 2, 583, 588, 3, 2, 2, 2, 584, 588, 
	5, 124, 63, 2, 585, 588, 5, 132, 67, 2, 586, 588, 5, 120, 61, 2, 587, 579, 
	3, 2, 2, 2, 587, 584, 3, 2, 2, 2, 587, 585, 3, 2, 2, 2, 587, 586, 3, 2, 
	2, 2, 588, 603, 3, 2, 2, 2, 589, 590, 12, 10, 2, 2, 590, 591, 7, 118, 2, 
	2, 591, 602, 5, 118, 60, 11, 592, 593, 12, 9, 2, 2, 593, 594, 7, 117, 2, 
	2, 594, 602, 5, 118, 60, 10, 595, 596, 12, 8, 2, 2, 596, 597, 7, 115, 2, 
	2, 597, 602, 5, 118, 60, 9, 598, 599, 12, 7, 2, 2, 599, 600, 7, 116, 2, 
	2, 600, 602, 5, 118, 60, 8, 601, 589, 3, 2, 2, 2, 601, 592, 3, 2, 2, 2, 
	601, 595, 3, 2, 2, 2, 601, 598, 3, 2, 2, 2, 602, 605, 3, 2, 2, 2, 603, 
	601, 3, 2, 2, 2, 603, 604, 3, 2, 2, 2, 604, 119, 3, 2, 2, 2, 605, 603, 
	3, 2, 2, 2, 606, 607, 5, 146, 74, 2, 607, 608, 5, 122, 62, 2, 608, 121, 
	3, 2, 2, 2, 609, 610, 9, 8, 2, 2, 610, 123, 3, 2, 2, 2, 611, 612, 5, 126, 
	64, 2, 612, 614, 7, 113, 2, 2, 613, 615, 5, 128, 65, 2, 614, 613, 3, 2, 
	2, 2, 614, 615, 3, 2, 2, 2, 615, 616, 3, 2, 2, 2, 616, 617, 7, 114, 2, 
	2, 617, 125, 3, 2, 2, 2, 618, 619, 9, 9, 2, 2, 619, 127, 3, 2, 2, 2, 620, 
	625, 5, 130, 66, 2, 621, 622, 7, 108, 2, 2, 622, 624, 5, 130, 66, 2, 623, 
	621, 3, 2, 2, 2, 624, 627, 3, 2, 2, 2, 625, 623, 3, 2, 2, 2, 625, 626, 
	3, 2, 2, 2, 626, 129, 3, 2, 2, 2, 627, 625, 3, 2, 2, 2, 628, 631, 5, 118, 
	60, 2, 629, 631, 5, 76, 39, 2, 630, 628, 3, 2, 2, 2, 630, 629, 3, 2, 2, 
	2, 631, 131, 3, 2, 2, 2, 632, 634, 5, 158, 80, 2, 633, 635, 5, 134, 68, 
	2, 634, 633, 3, 2, 2, 2, 634, 635, 3, 2, 2, 2, 635, 639, 3, 2, 2, 2, 636, 
	639, 5, 148, 75, 2, 637, 639, 5, 146, 74, 2, 638, 632, 3, 2, 2, 2, 638, 
	636, 3, 2, 2, 2, 638, 637, 3, 2, 2, 2, 639, 133, 3, 2, 2, 2, 640, 641, 
	7, 111, 2, 2, 641, 642, 5, 76, 39, 2, 642, 643, 7, 112, 2, 2, 643, 135, 
	3, 2, 2, 2, 644, 645, 5, 144, 73, 2, 645, 137, 3, 2, 2, 2, 646, 647, 7, 
	109, 2, 2, 647, 652, 5, 140, 71, 2, 648, 649, 7, 108, 2, 2, 649, 651, 5, 
	140, 71, 2, 650, 648, 3, 2, 2, 2, 651, 654, 3, 2, 2, 2, 652, 650, 3, 2, 
	2, 2, 652, 653, 3, 2, 2, 2, 653, 655, 3, 2, 2, 2, 654, 652, 3, 2, 2, 2, 
	655, 656, 7, 110, 2, 2, 656, 660, 3, 2, 2, 2, 657, 658, 7, 109, 2, 2, 658, 
	660, 7, 110, 2, 2, 659, 646, 3, 2, 2, 2, 659, 657, 3, 2, 2, 2, 660, 139, 
	3, 2, 2, 2, 661, 662, 7, 6, 2, 2, 662, 663, 7, 98, 2, 2, 663, 664, 5, 144, 
	73, 2, 664, 141, 3, 2, 2, 2, 665, 666, 7, 111, 2, 2, 666, 671, 5, 144, 
	73, 2, 667, 668, 7, 108, 2, 2, 668, 670, 5, 144, 73, 2, 669, 667, 3, 2, 
	2, 2, 670, 673, 3, 2, 2, 2, 671, 669, 3, 2, 2, 2, 671, 672, 3, 2, 2, 2, 
	672, 674, 3, 2, 2, 2, 673, 671, 3, 2, 2, 2, 674, 675, 7, 112, 2, 2, 675, 
	679, 3, 2, 2, 2, 676, 677, 7, 111, 2, 2, 677, 679, 7, 112, 2, 2, 678, 665, 
	3, 2, 2, 2, 678, 676, 3, 2, 2, 2, 679, 143, 3, 2, 2, 2, 680, 689, 7, 6, 
	2, 2, 681, 689, 5, 146, 74, 2, 682, 689, 5, 148, 75, 2, 683, 689, 5, 138, 
	70, 2, 684, 689, 5, 142, 72, 2, 685, 689, 7, 3, 2, 2, 686, 689, 7, 4, 2, 
	2, 687, 689, 7, 5, 2, 2, 688, 680, 3, 2, 2, 2, 688, 681, 3, 2, 2, 2, 688, 
	682, 3, 2, 2, 2, 688, 683, 3, 2, 2, 2, 688, 684, 3, 2, 2, 2, 688, 685, 
	3, 2, 2, 2, 688, 686, 3, 2, 2, 2, 688, 687, 3, 2, 2, 2, 689, 145, 3, 2, 
	2, 2, 690, 692, 9, 10, 2, 2, 691, 690, 3, 2, 2, 2, 691, 692, 3, 2, 2, 2, 
	692, 693, 3, 2, 2, 2, 693, 694, 7, 122, 2, 2, 694, 147, 3, 2, 2, 2, 695, 
	697, 9, 10, 2, 2, 696, 695, 3, 2, 2, 2, 696, 697, 3, 2, 2, 2, 697, 698, 
	3, 2, 2, 2, 698, 699, 7, 123, 2, 2, 699, 149, 3, 2, 2, 2, 700, 701, 7, 
	53, 2, 2, 701, 702, 7, 122, 2, 2, 702, 151, 3, 2, 2, 2, 703, 704, 5, 158, 
	80, 2, 704, 153, 3, 2, 2, 2, 705, 706, 5, 158, 80, 2, 706, 155, 3, 2, 2, 
	2, 707, 708, 5, 158, 80, 2, 708, 157, 3, 2, 2, 2, 709, 712, 7, 121, 2, 
	2, 710, 712, 5, 160, 81, 2, 711, 709, 3, 2, 2, 2, 711, 710, 3, 2, 2, 2, 
	712, 720, 3, 2, 2, 2, 713, 716, 7, 97, 2, 2, 714, 717, 7, 121, 2, 2, 715, 
	717, 5, 160, 81, 2, 716, 714, 3, 2, 2, 2, 716, 715, 3, 2, 2, 2, 717, 719, 
	3, 2, 2, 2, 718, 713, 3, 2, 2, 2, 719, 722, 3, 2, 2, 2, 720, 718, 3, 2, 
	2, 2, 720, 721, 3, 2, 2, 2, 721, 159, 3, 2, 2, 2, 722, 720, 3, 2, 2, 2, 
	723, 724, 9, 11, 2, 2, 724, 161, 3, 2, 2, 2, 64, 186, 225, 230, 241, 246, 
	260, 265, 287, 290, 296, 302, 305, 325, 328, 339, 344, 347, 350, 353, 356, 
	366, 371, 392, 405, 407, 423, 431, 437, 444, 452, 466, 472, 478, 482, 487, 
	499, 502, 509, 518, 530, 538, 550, 558, 577, 587, 601, 603, 614, 625, 630, 
	634, 638, 652, 659, 671, 678, 688, 691, 696, 711, 716, 720,
}
var literalNames = []string{
	"", "'true'", "'false'", "'null'", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
	"", "", "", "", "'m'", "", "", "", "'M'", "", "'.'", "':'", "'='", "'<>'", 
	"'!='", "'>'", "'>='", "'<'", "'<='", "'=~'", "'!~'", "','", "'{'", "'}'", 
	"'['", "']'", "'('", "')'", "'+'", "'-'", "'/'", "'*'", "'%'", "'_'",
}
var symbolicNames = []string{
	"", "", "", "", "STRING", "WS", "T_CREATE", "T_UPDATE", "T_SET", "T_DROP", 
	"T_INTERVAL", "T_INTERVAL_NAME", "T_SHARD", "T_REPLICATION", "T_TTL", "T_META_TTL", 
	"T_PAST_TTL", "T_FUTURE_TTL", "T_KILL", "T_ON", "T_SHOW", "T_USE", "T_STATE_REPO", 
	"T_STATE_MACHINE", "T_MASTER", "T_METADATA", "T_TYPES", "T_TYPE", "T_STORAGES", 
	"T_STORAGE", "T_BROKER", "T_ALIVE", "T_SCHEMAS", "T_DATASBAE", "T_DATASBAES", 
	"T_NAMESPACE", "T_NAMESPACES", "T_NODE", "T_METRICS", "T_METRIC", "T_FIELD", 
	"T_FIELDS", "T_TAG", "T_INFO", "T_KEYS", "T_KEY", "T_WITH", "T_VALUES", 
	"T_VALUE", "T_FROM", "T_WHERE", "T_LIMIT", "T_QUERIES", "T_QUERY", "T_EXPLAIN", 
	"T_WITH_VALUE", "T_SELECT", "T_AS", "T_AND", "T_OR", "T_FILL", "T_NULL", 
	"T_PREVIOUS", "T_ORDER", "T_ASC", "T_DESC", "T_LIKE", "T_NOT", "T_BETWEEN", 
	"T_IS", "T_GROUP", "T_HAVING", "T_BY", "T_FOR", "T_STATS", "T_TIME", "T_NOW", 
	"T_IN", "T_LOG", "T_PROFILE", "T_SUM", "T_MIN", "T_MAX", "T_COUNT", "T_AVG", 
	"T_STDDEV", "T_QUANTILE", "T_RATE", "T_SECOND", "T_MINUTE", "T_HOUR", "T_DAY", 
	"T_WEEK", "T_MONTH", "T_YEAR", "T_DOT", "T_COLON", "T_EQUAL", "T_NOTEQUAL", 
	"T_NOTEQUAL2", "T_GREATER", "T_GREATEREQUAL", "T_LESS", "T_LESSEQUAL", 
	"T_REGEXP", "T_NEQREGEXP", "T_COMMA", "T_OPEN_B", "T_CLOSE_B", "T_OPEN_SB", 
	"T_CLOSE_SB", "T_OPEN_P", "T_CLOSE_P", "T_ADD", "T_SUB", "T_DIV", "T_MUL", 
	"T_MOD", "T_UNDERLINE", "L_ID", "L_INT", "L_DEC",
}

var ruleNames = []string{
	"statement", "statementList", "useStmt", "showMasterStmt", "showStoragesStmt", 
	"showMetadataTypesStmt", "showBrokerMetaStmt", "showMasterMetaStmt", "showStorageMetaStmt", 
	"showAliveStmt", "showReplicationStmt", "showBrokerMetricStmt", "showStorageMetricStmt", 
	"createStorageStmt", "showSchemasStmt", "createDatabaseStmt", "showDatabaseStmt", 
	"showNameSpacesStmt", "showMetricsStmt", "showFieldsStmt", "showTagKeysStmt", 
	"showTagValuesStmt", "prefix", "withTagKey", "namespace", "source", "queryStmt", 
	"selectExpr", "fields", "field", "alias", "storageFilter", "databaseFilter", 
	"typeFilter", "fromClause", "whereClause", "conditionExpr", "tagFilterExpr", 
	"tagValueList", "metricListFilter", "metricList", "timeRangeExpr", "timeExpr", 
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
	SQLParserT_STATE_REPO = 22
	SQLParserT_STATE_MACHINE = 23
	SQLParserT_MASTER = 24
	SQLParserT_METADATA = 25
	SQLParserT_TYPES = 26
	SQLParserT_TYPE = 27
	SQLParserT_STORAGES = 28
	SQLParserT_STORAGE = 29
	SQLParserT_BROKER = 30
	SQLParserT_ALIVE = 31
	SQLParserT_SCHEMAS = 32
	SQLParserT_DATASBAE = 33
	SQLParserT_DATASBAES = 34
	SQLParserT_NAMESPACE = 35
	SQLParserT_NAMESPACES = 36
	SQLParserT_NODE = 37
	SQLParserT_METRICS = 38
	SQLParserT_METRIC = 39
	SQLParserT_FIELD = 40
	SQLParserT_FIELDS = 41
	SQLParserT_TAG = 42
	SQLParserT_INFO = 43
	SQLParserT_KEYS = 44
	SQLParserT_KEY = 45
	SQLParserT_WITH = 46
	SQLParserT_VALUES = 47
	SQLParserT_VALUE = 48
	SQLParserT_FROM = 49
	SQLParserT_WHERE = 50
	SQLParserT_LIMIT = 51
	SQLParserT_QUERIES = 52
	SQLParserT_QUERY = 53
	SQLParserT_EXPLAIN = 54
	SQLParserT_WITH_VALUE = 55
	SQLParserT_SELECT = 56
	SQLParserT_AS = 57
	SQLParserT_AND = 58
	SQLParserT_OR = 59
	SQLParserT_FILL = 60
	SQLParserT_NULL = 61
	SQLParserT_PREVIOUS = 62
	SQLParserT_ORDER = 63
	SQLParserT_ASC = 64
	SQLParserT_DESC = 65
	SQLParserT_LIKE = 66
	SQLParserT_NOT = 67
	SQLParserT_BETWEEN = 68
	SQLParserT_IS = 69
	SQLParserT_GROUP = 70
	SQLParserT_HAVING = 71
	SQLParserT_BY = 72
	SQLParserT_FOR = 73
	SQLParserT_STATS = 74
	SQLParserT_TIME = 75
	SQLParserT_NOW = 76
	SQLParserT_IN = 77
	SQLParserT_LOG = 78
	SQLParserT_PROFILE = 79
	SQLParserT_SUM = 80
	SQLParserT_MIN = 81
	SQLParserT_MAX = 82
	SQLParserT_COUNT = 83
	SQLParserT_AVG = 84
	SQLParserT_STDDEV = 85
	SQLParserT_QUANTILE = 86
	SQLParserT_RATE = 87
	SQLParserT_SECOND = 88
	SQLParserT_MINUTE = 89
	SQLParserT_HOUR = 90
	SQLParserT_DAY = 91
	SQLParserT_WEEK = 92
	SQLParserT_MONTH = 93
	SQLParserT_YEAR = 94
	SQLParserT_DOT = 95
	SQLParserT_COLON = 96
	SQLParserT_EQUAL = 97
	SQLParserT_NOTEQUAL = 98
	SQLParserT_NOTEQUAL2 = 99
	SQLParserT_GREATER = 100
	SQLParserT_GREATEREQUAL = 101
	SQLParserT_LESS = 102
	SQLParserT_LESSEQUAL = 103
	SQLParserT_REGEXP = 104
	SQLParserT_NEQREGEXP = 105
	SQLParserT_COMMA = 106
	SQLParserT_OPEN_B = 107
	SQLParserT_CLOSE_B = 108
	SQLParserT_OPEN_SB = 109
	SQLParserT_CLOSE_SB = 110
	SQLParserT_OPEN_P = 111
	SQLParserT_CLOSE_P = 112
	SQLParserT_ADD = 113
	SQLParserT_SUB = 114
	SQLParserT_DIV = 115
	SQLParserT_MUL = 116
	SQLParserT_MOD = 117
	SQLParserT_UNDERLINE = 118
	SQLParserL_ID = 119
	SQLParserL_INT = 120
	SQLParserL_DEC = 121
)

// SQLParser rules.
const (
	SQLParserRULE_statement = 0
	SQLParserRULE_statementList = 1
	SQLParserRULE_useStmt = 2
	SQLParserRULE_showMasterStmt = 3
	SQLParserRULE_showStoragesStmt = 4
	SQLParserRULE_showMetadataTypesStmt = 5
	SQLParserRULE_showBrokerMetaStmt = 6
	SQLParserRULE_showMasterMetaStmt = 7
	SQLParserRULE_showStorageMetaStmt = 8
	SQLParserRULE_showAliveStmt = 9
	SQLParserRULE_showReplicationStmt = 10
	SQLParserRULE_showBrokerMetricStmt = 11
	SQLParserRULE_showStorageMetricStmt = 12
	SQLParserRULE_createStorageStmt = 13
	SQLParserRULE_showSchemasStmt = 14
	SQLParserRULE_createDatabaseStmt = 15
	SQLParserRULE_showDatabaseStmt = 16
	SQLParserRULE_showNameSpacesStmt = 17
	SQLParserRULE_showMetricsStmt = 18
	SQLParserRULE_showFieldsStmt = 19
	SQLParserRULE_showTagKeysStmt = 20
	SQLParserRULE_showTagValuesStmt = 21
	SQLParserRULE_prefix = 22
	SQLParserRULE_withTagKey = 23
	SQLParserRULE_namespace = 24
	SQLParserRULE_source = 25
	SQLParserRULE_queryStmt = 26
	SQLParserRULE_selectExpr = 27
	SQLParserRULE_fields = 28
	SQLParserRULE_field = 29
	SQLParserRULE_alias = 30
	SQLParserRULE_storageFilter = 31
	SQLParserRULE_databaseFilter = 32
	SQLParserRULE_typeFilter = 33
	SQLParserRULE_fromClause = 34
	SQLParserRULE_whereClause = 35
	SQLParserRULE_conditionExpr = 36
	SQLParserRULE_tagFilterExpr = 37
	SQLParserRULE_tagValueList = 38
	SQLParserRULE_metricListFilter = 39
	SQLParserRULE_metricList = 40
	SQLParserRULE_timeRangeExpr = 41
	SQLParserRULE_timeExpr = 42
	SQLParserRULE_nowExpr = 43
	SQLParserRULE_nowFunc = 44
	SQLParserRULE_groupByClause = 45
	SQLParserRULE_groupByKeys = 46
	SQLParserRULE_groupByKey = 47
	SQLParserRULE_fillOption = 48
	SQLParserRULE_orderByClause = 49
	SQLParserRULE_sortField = 50
	SQLParserRULE_sortFields = 51
	SQLParserRULE_havingClause = 52
	SQLParserRULE_boolExpr = 53
	SQLParserRULE_boolExprLogicalOp = 54
	SQLParserRULE_boolExprAtom = 55
	SQLParserRULE_binaryExpr = 56
	SQLParserRULE_binaryOperator = 57
	SQLParserRULE_fieldExpr = 58
	SQLParserRULE_durationLit = 59
	SQLParserRULE_intervalItem = 60
	SQLParserRULE_exprFunc = 61
	SQLParserRULE_funcName = 62
	SQLParserRULE_exprFuncParams = 63
	SQLParserRULE_funcParam = 64
	SQLParserRULE_exprAtom = 65
	SQLParserRULE_identFilter = 66
	SQLParserRULE_json = 67
	SQLParserRULE_obj = 68
	SQLParserRULE_pair = 69
	SQLParserRULE_arr = 70
	SQLParserRULE_value = 71
	SQLParserRULE_intNumber = 72
	SQLParserRULE_decNumber = 73
	SQLParserRULE_limitClause = 74
	SQLParserRULE_metricName = 75
	SQLParserRULE_tagKey = 76
	SQLParserRULE_tagValue = 77
	SQLParserRULE_ident = 78
	SQLParserRULE_nonReservedWords = 79
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
		p.SetState(160)
		p.StatementList()
	}
	{
		p.SetState(161)
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

func (s *StatementListContext) ShowMetadataTypesStmt() IShowMetadataTypesStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowMetadataTypesStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowMetadataTypesStmtContext)
}

func (s *StatementListContext) ShowBrokerMetaStmt() IShowBrokerMetaStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowBrokerMetaStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowBrokerMetaStmtContext)
}

func (s *StatementListContext) ShowMasterMetaStmt() IShowMasterMetaStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowMasterMetaStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowMasterMetaStmtContext)
}

func (s *StatementListContext) ShowStorageMetaStmt() IShowStorageMetaStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowStorageMetaStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowStorageMetaStmtContext)
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

func (s *StatementListContext) ShowBrokerMetricStmt() IShowBrokerMetricStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowBrokerMetricStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowBrokerMetricStmtContext)
}

func (s *StatementListContext) ShowStorageMetricStmt() IShowStorageMetricStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowStorageMetricStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowStorageMetricStmtContext)
}

func (s *StatementListContext) CreateStorageStmt() ICreateStorageStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICreateStorageStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICreateStorageStmtContext)
}

func (s *StatementListContext) ShowReplicationStmt() IShowReplicationStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowReplicationStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowReplicationStmtContext)
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

	p.SetState(184)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(163)
			p.ShowMasterStmt()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(164)
			p.ShowMetadataTypesStmt()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(165)
			p.ShowBrokerMetaStmt()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(166)
			p.ShowMasterMetaStmt()
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(167)
			p.ShowStorageMetaStmt()
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(168)
			p.ShowStoragesStmt()
		}


	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(169)
			p.ShowAliveStmt()
		}


	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(170)
			p.ShowBrokerMetricStmt()
		}


	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(171)
			p.ShowStorageMetricStmt()
		}


	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(172)
			p.CreateStorageStmt()
		}


	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(173)
			p.ShowReplicationStmt()
		}


	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(174)
			p.ShowSchemasStmt()
		}


	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(175)
			p.ShowDatabaseStmt()
		}


	case 14:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(176)
			p.UseStmt()
		}


	case 15:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(177)
			p.ShowNameSpacesStmt()
		}


	case 16:
		p.EnterOuterAlt(localctx, 16)
		{
			p.SetState(178)
			p.ShowMetricsStmt()
		}


	case 17:
		p.EnterOuterAlt(localctx, 17)
		{
			p.SetState(179)
			p.ShowFieldsStmt()
		}


	case 18:
		p.EnterOuterAlt(localctx, 18)
		{
			p.SetState(180)
			p.ShowTagKeysStmt()
		}


	case 19:
		p.EnterOuterAlt(localctx, 19)
		{
			p.SetState(181)
			p.ShowTagValuesStmt()
		}


	case 20:
		p.EnterOuterAlt(localctx, 20)
		{
			p.SetState(182)
			p.QueryStmt()
		}


	case 21:
		p.EnterOuterAlt(localctx, 21)
		{
			p.SetState(183)
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
		p.SetState(186)
		p.Match(SQLParserT_USE)
	}
	{
		p.SetState(187)
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
		p.SetState(189)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(190)
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
		p.SetState(192)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(193)
		p.Match(SQLParserT_STORAGES)
	}



	return localctx
}


// IShowMetadataTypesStmtContext is an interface to support dynamic dispatch.
type IShowMetadataTypesStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowMetadataTypesStmtContext differentiates from other interfaces.
	IsShowMetadataTypesStmtContext()
}

type ShowMetadataTypesStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowMetadataTypesStmtContext() *ShowMetadataTypesStmtContext {
	var p = new(ShowMetadataTypesStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showMetadataTypesStmt
	return p
}

func (*ShowMetadataTypesStmtContext) IsShowMetadataTypesStmtContext() {}

func NewShowMetadataTypesStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowMetadataTypesStmtContext {
	var p = new(ShowMetadataTypesStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showMetadataTypesStmt

	return p
}

func (s *ShowMetadataTypesStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowMetadataTypesStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowMetadataTypesStmtContext) T_METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METADATA, 0)
}

func (s *ShowMetadataTypesStmtContext) T_TYPES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TYPES, 0)
}

func (s *ShowMetadataTypesStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMetadataTypesStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowMetadataTypesStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowMetadataTypesStmt(s)
	}
}

func (s *ShowMetadataTypesStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowMetadataTypesStmt(s)
	}
}




func (p *SQLParser) ShowMetadataTypesStmt() (localctx IShowMetadataTypesStmtContext) {
	localctx = NewShowMetadataTypesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SQLParserRULE_showMetadataTypesStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(195)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(196)
		p.Match(SQLParserT_METADATA)
	}
	{
		p.SetState(197)
		p.Match(SQLParserT_TYPES)
	}



	return localctx
}


// IShowBrokerMetaStmtContext is an interface to support dynamic dispatch.
type IShowBrokerMetaStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowBrokerMetaStmtContext differentiates from other interfaces.
	IsShowBrokerMetaStmtContext()
}

type ShowBrokerMetaStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowBrokerMetaStmtContext() *ShowBrokerMetaStmtContext {
	var p = new(ShowBrokerMetaStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showBrokerMetaStmt
	return p
}

func (*ShowBrokerMetaStmtContext) IsShowBrokerMetaStmtContext() {}

func NewShowBrokerMetaStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowBrokerMetaStmtContext {
	var p = new(ShowBrokerMetaStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showBrokerMetaStmt

	return p
}

func (s *ShowBrokerMetaStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowBrokerMetaStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowBrokerMetaStmtContext) T_BROKER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BROKER, 0)
}

func (s *ShowBrokerMetaStmtContext) T_METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METADATA, 0)
}

func (s *ShowBrokerMetaStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowBrokerMetaStmtContext) Source() ISourceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISourceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISourceContext)
}

func (s *ShowBrokerMetaStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowBrokerMetaStmtContext) TypeFilter() ITypeFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeFilterContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeFilterContext)
}

func (s *ShowBrokerMetaStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowBrokerMetaStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowBrokerMetaStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowBrokerMetaStmt(s)
	}
}

func (s *ShowBrokerMetaStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowBrokerMetaStmt(s)
	}
}




func (p *SQLParser) ShowBrokerMetaStmt() (localctx IShowBrokerMetaStmtContext) {
	localctx = NewShowBrokerMetaStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SQLParserRULE_showBrokerMetaStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(199)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(200)
		p.Match(SQLParserT_BROKER)
	}
	{
		p.SetState(201)
		p.Match(SQLParserT_METADATA)
	}
	{
		p.SetState(202)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(203)
		p.Source()
	}
	{
		p.SetState(204)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(205)
		p.TypeFilter()
	}



	return localctx
}


// IShowMasterMetaStmtContext is an interface to support dynamic dispatch.
type IShowMasterMetaStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowMasterMetaStmtContext differentiates from other interfaces.
	IsShowMasterMetaStmtContext()
}

type ShowMasterMetaStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowMasterMetaStmtContext() *ShowMasterMetaStmtContext {
	var p = new(ShowMasterMetaStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showMasterMetaStmt
	return p
}

func (*ShowMasterMetaStmtContext) IsShowMasterMetaStmtContext() {}

func NewShowMasterMetaStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowMasterMetaStmtContext {
	var p = new(ShowMasterMetaStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showMasterMetaStmt

	return p
}

func (s *ShowMasterMetaStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowMasterMetaStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowMasterMetaStmtContext) T_MASTER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MASTER, 0)
}

func (s *ShowMasterMetaStmtContext) T_METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METADATA, 0)
}

func (s *ShowMasterMetaStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowMasterMetaStmtContext) Source() ISourceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISourceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISourceContext)
}

func (s *ShowMasterMetaStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowMasterMetaStmtContext) TypeFilter() ITypeFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeFilterContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeFilterContext)
}

func (s *ShowMasterMetaStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMasterMetaStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowMasterMetaStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowMasterMetaStmt(s)
	}
}

func (s *ShowMasterMetaStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowMasterMetaStmt(s)
	}
}




func (p *SQLParser) ShowMasterMetaStmt() (localctx IShowMasterMetaStmtContext) {
	localctx = NewShowMasterMetaStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SQLParserRULE_showMasterMetaStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(207)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(208)
		p.Match(SQLParserT_MASTER)
	}
	{
		p.SetState(209)
		p.Match(SQLParserT_METADATA)
	}
	{
		p.SetState(210)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(211)
		p.Source()
	}
	{
		p.SetState(212)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(213)
		p.TypeFilter()
	}



	return localctx
}


// IShowStorageMetaStmtContext is an interface to support dynamic dispatch.
type IShowStorageMetaStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowStorageMetaStmtContext differentiates from other interfaces.
	IsShowStorageMetaStmtContext()
}

type ShowStorageMetaStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowStorageMetaStmtContext() *ShowStorageMetaStmtContext {
	var p = new(ShowStorageMetaStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showStorageMetaStmt
	return p
}

func (*ShowStorageMetaStmtContext) IsShowStorageMetaStmtContext() {}

func NewShowStorageMetaStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowStorageMetaStmtContext {
	var p = new(ShowStorageMetaStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showStorageMetaStmt

	return p
}

func (s *ShowStorageMetaStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowStorageMetaStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowStorageMetaStmtContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *ShowStorageMetaStmtContext) T_METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METADATA, 0)
}

func (s *ShowStorageMetaStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowStorageMetaStmtContext) Source() ISourceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISourceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISourceContext)
}

func (s *ShowStorageMetaStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowStorageMetaStmtContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *ShowStorageMetaStmtContext) AllStorageFilter() []IStorageFilterContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStorageFilterContext)(nil)).Elem())
	var tst = make([]IStorageFilterContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStorageFilterContext)
		}
	}

	return tst
}

func (s *ShowStorageMetaStmtContext) StorageFilter(i int) IStorageFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStorageFilterContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStorageFilterContext)
}

func (s *ShowStorageMetaStmtContext) AllTypeFilter() []ITypeFilterContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITypeFilterContext)(nil)).Elem())
	var tst = make([]ITypeFilterContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITypeFilterContext)
		}
	}

	return tst
}

func (s *ShowStorageMetaStmtContext) TypeFilter(i int) ITypeFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeFilterContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITypeFilterContext)
}

func (s *ShowStorageMetaStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStorageMetaStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowStorageMetaStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowStorageMetaStmt(s)
	}
}

func (s *ShowStorageMetaStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowStorageMetaStmt(s)
	}
}




func (p *SQLParser) ShowStorageMetaStmt() (localctx IShowStorageMetaStmtContext) {
	localctx = NewShowStorageMetaStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SQLParserRULE_showStorageMetaStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(215)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(216)
		p.Match(SQLParserT_STORAGE)
	}
	{
		p.SetState(217)
		p.Match(SQLParserT_METADATA)
	}
	{
		p.SetState(218)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(219)
		p.Source()
	}
	{
		p.SetState(220)
		p.Match(SQLParserT_WHERE)
	}
	p.SetState(223)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(221)
			p.StorageFilter()
		}


	case SQLParserT_TYPE:
		{
			p.SetState(222)
			p.TypeFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(225)
		p.Match(SQLParserT_AND)
	}
	p.SetState(228)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(226)
			p.StorageFilter()
		}


	case SQLParserT_TYPE:
		{
			p.SetState(227)
			p.TypeFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 18, SQLParserRULE_showAliveStmt)
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
		p.SetState(230)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(231)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserT_STORAGE || _la == SQLParserT_BROKER) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(232)
		p.Match(SQLParserT_ALIVE)
	}



	return localctx
}


// IShowReplicationStmtContext is an interface to support dynamic dispatch.
type IShowReplicationStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowReplicationStmtContext differentiates from other interfaces.
	IsShowReplicationStmtContext()
}

type ShowReplicationStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowReplicationStmtContext() *ShowReplicationStmtContext {
	var p = new(ShowReplicationStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showReplicationStmt
	return p
}

func (*ShowReplicationStmtContext) IsShowReplicationStmtContext() {}

func NewShowReplicationStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowReplicationStmtContext {
	var p = new(ShowReplicationStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showReplicationStmt

	return p
}

func (s *ShowReplicationStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowReplicationStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowReplicationStmtContext) T_REPLICATION() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REPLICATION, 0)
}

func (s *ShowReplicationStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowReplicationStmtContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *ShowReplicationStmtContext) AllStorageFilter() []IStorageFilterContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStorageFilterContext)(nil)).Elem())
	var tst = make([]IStorageFilterContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStorageFilterContext)
		}
	}

	return tst
}

func (s *ShowReplicationStmtContext) StorageFilter(i int) IStorageFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStorageFilterContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStorageFilterContext)
}

func (s *ShowReplicationStmtContext) AllDatabaseFilter() []IDatabaseFilterContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IDatabaseFilterContext)(nil)).Elem())
	var tst = make([]IDatabaseFilterContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IDatabaseFilterContext)
		}
	}

	return tst
}

func (s *ShowReplicationStmtContext) DatabaseFilter(i int) IDatabaseFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDatabaseFilterContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IDatabaseFilterContext)
}

func (s *ShowReplicationStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowReplicationStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowReplicationStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowReplicationStmt(s)
	}
}

func (s *ShowReplicationStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowReplicationStmt(s)
	}
}




func (p *SQLParser) ShowReplicationStmt() (localctx IShowReplicationStmtContext) {
	localctx = NewShowReplicationStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SQLParserRULE_showReplicationStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(234)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(235)
		p.Match(SQLParserT_REPLICATION)
	}
	{
		p.SetState(236)
		p.Match(SQLParserT_WHERE)
	}
	p.SetState(239)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(237)
			p.StorageFilter()
		}


	case SQLParserT_DATASBAE:
		{
			p.SetState(238)
			p.DatabaseFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(241)
		p.Match(SQLParserT_AND)
	}
	p.SetState(244)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(242)
			p.StorageFilter()
		}


	case SQLParserT_DATASBAE:
		{
			p.SetState(243)
			p.DatabaseFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}



	return localctx
}


// IShowBrokerMetricStmtContext is an interface to support dynamic dispatch.
type IShowBrokerMetricStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowBrokerMetricStmtContext differentiates from other interfaces.
	IsShowBrokerMetricStmtContext()
}

type ShowBrokerMetricStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowBrokerMetricStmtContext() *ShowBrokerMetricStmtContext {
	var p = new(ShowBrokerMetricStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showBrokerMetricStmt
	return p
}

func (*ShowBrokerMetricStmtContext) IsShowBrokerMetricStmtContext() {}

func NewShowBrokerMetricStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowBrokerMetricStmtContext {
	var p = new(ShowBrokerMetricStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showBrokerMetricStmt

	return p
}

func (s *ShowBrokerMetricStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowBrokerMetricStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowBrokerMetricStmtContext) T_BROKER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BROKER, 0)
}

func (s *ShowBrokerMetricStmtContext) T_METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRIC, 0)
}

func (s *ShowBrokerMetricStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowBrokerMetricStmtContext) MetricListFilter() IMetricListFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetricListFilterContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetricListFilterContext)
}

func (s *ShowBrokerMetricStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowBrokerMetricStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowBrokerMetricStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowBrokerMetricStmt(s)
	}
}

func (s *ShowBrokerMetricStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowBrokerMetricStmt(s)
	}
}




func (p *SQLParser) ShowBrokerMetricStmt() (localctx IShowBrokerMetricStmtContext) {
	localctx = NewShowBrokerMetricStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SQLParserRULE_showBrokerMetricStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
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
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(247)
		p.Match(SQLParserT_BROKER)
	}
	{
		p.SetState(248)
		p.Match(SQLParserT_METRIC)
	}
	{
		p.SetState(249)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(250)
		p.MetricListFilter()
	}



	return localctx
}


// IShowStorageMetricStmtContext is an interface to support dynamic dispatch.
type IShowStorageMetricStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowStorageMetricStmtContext differentiates from other interfaces.
	IsShowStorageMetricStmtContext()
}

type ShowStorageMetricStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowStorageMetricStmtContext() *ShowStorageMetricStmtContext {
	var p = new(ShowStorageMetricStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showStorageMetricStmt
	return p
}

func (*ShowStorageMetricStmtContext) IsShowStorageMetricStmtContext() {}

func NewShowStorageMetricStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowStorageMetricStmtContext {
	var p = new(ShowStorageMetricStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showStorageMetricStmt

	return p
}

func (s *ShowStorageMetricStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowStorageMetricStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowStorageMetricStmtContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *ShowStorageMetricStmtContext) T_METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRIC, 0)
}

func (s *ShowStorageMetricStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowStorageMetricStmtContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *ShowStorageMetricStmtContext) AllStorageFilter() []IStorageFilterContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStorageFilterContext)(nil)).Elem())
	var tst = make([]IStorageFilterContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStorageFilterContext)
		}
	}

	return tst
}

func (s *ShowStorageMetricStmtContext) StorageFilter(i int) IStorageFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStorageFilterContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStorageFilterContext)
}

func (s *ShowStorageMetricStmtContext) AllMetricListFilter() []IMetricListFilterContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IMetricListFilterContext)(nil)).Elem())
	var tst = make([]IMetricListFilterContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IMetricListFilterContext)
		}
	}

	return tst
}

func (s *ShowStorageMetricStmtContext) MetricListFilter(i int) IMetricListFilterContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetricListFilterContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IMetricListFilterContext)
}

func (s *ShowStorageMetricStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStorageMetricStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowStorageMetricStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowStorageMetricStmt(s)
	}
}

func (s *ShowStorageMetricStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowStorageMetricStmt(s)
	}
}




func (p *SQLParser) ShowStorageMetricStmt() (localctx IShowStorageMetricStmtContext) {
	localctx = NewShowStorageMetricStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, SQLParserRULE_showStorageMetricStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
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
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(253)
		p.Match(SQLParserT_STORAGE)
	}
	{
		p.SetState(254)
		p.Match(SQLParserT_METRIC)
	}
	{
		p.SetState(255)
		p.Match(SQLParserT_WHERE)
	}
	p.SetState(258)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(256)
			p.StorageFilter()
		}


	case SQLParserT_METRIC:
		{
			p.SetState(257)
			p.MetricListFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(260)
		p.Match(SQLParserT_AND)
	}
	p.SetState(263)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(261)
			p.StorageFilter()
		}


	case SQLParserT_METRIC:
		{
			p.SetState(262)
			p.MetricListFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 26, SQLParserRULE_createStorageStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(265)
		p.Match(SQLParserT_CREATE)
	}
	{
		p.SetState(266)
		p.Match(SQLParserT_STORAGE)
	}
	{
		p.SetState(267)
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
	p.EnterRule(localctx, 28, SQLParserRULE_showSchemasStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(269)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(270)
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
	p.EnterRule(localctx, 30, SQLParserRULE_createDatabaseStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(272)
		p.Match(SQLParserT_CREATE)
	}
	{
		p.SetState(273)
		p.Match(SQLParserT_DATASBAE)
	}
	{
		p.SetState(274)
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
	p.EnterRule(localctx, 32, SQLParserRULE_showDatabaseStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(276)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(277)
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
	p.EnterRule(localctx, 34, SQLParserRULE_showNameSpacesStmt)
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
		p.SetState(279)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(280)
		p.Match(SQLParserT_NAMESPACES)
	}
	p.SetState(285)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(281)
			p.Match(SQLParserT_WHERE)
		}
		{
			p.SetState(282)
			p.Match(SQLParserT_NAMESPACE)
		}
		{
			p.SetState(283)
			p.Match(SQLParserT_EQUAL)
		}
		{
			p.SetState(284)
			p.Prefix()
		}

	}
	p.SetState(288)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(287)
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
	p.EnterRule(localctx, 36, SQLParserRULE_showMetricsStmt)
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
		p.SetState(290)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(291)
		p.Match(SQLParserT_METRICS)
	}
	p.SetState(294)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(292)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(293)
			p.Namespace()
		}

	}
	p.SetState(300)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(296)
			p.Match(SQLParserT_WHERE)
		}
		{
			p.SetState(297)
			p.Match(SQLParserT_METRIC)
		}
		{
			p.SetState(298)
			p.Match(SQLParserT_EQUAL)
		}
		{
			p.SetState(299)
			p.Prefix()
		}

	}
	p.SetState(303)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(302)
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
	p.EnterRule(localctx, 38, SQLParserRULE_showFieldsStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(305)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(306)
		p.Match(SQLParserT_FIELDS)
	}
	{
		p.SetState(307)
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
	p.EnterRule(localctx, 40, SQLParserRULE_showTagKeysStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
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
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(310)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(311)
		p.Match(SQLParserT_KEYS)
	}
	{
		p.SetState(312)
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
	p.EnterRule(localctx, 42, SQLParserRULE_showTagValuesStmt)
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
		p.SetState(314)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(315)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(316)
		p.Match(SQLParserT_VALUES)
	}
	{
		p.SetState(317)
		p.FromClause()
	}
	{
		p.SetState(318)
		p.Match(SQLParserT_WITH)
	}
	{
		p.SetState(319)
		p.Match(SQLParserT_KEY)
	}
	{
		p.SetState(320)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(321)
		p.WithTagKey()
	}
	p.SetState(323)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(322)
			p.WhereClause()
		}

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
	p.EnterRule(localctx, 44, SQLParserRULE_prefix)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
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
	p.EnterRule(localctx, 46, SQLParserRULE_withTagKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(330)
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
	p.EnterRule(localctx, 48, SQLParserRULE_namespace)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(332)
		p.Ident()
	}



	return localctx
}


// ISourceContext is an interface to support dynamic dispatch.
type ISourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSourceContext differentiates from other interfaces.
	IsSourceContext()
}

type SourceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySourceContext() *SourceContext {
	var p = new(SourceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_source
	return p
}

func (*SourceContext) IsSourceContext() {}

func NewSourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SourceContext {
	var p = new(SourceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_source

	return p
}

func (s *SourceContext) GetParser() antlr.Parser { return s.parser }

func (s *SourceContext) T_STATE_MACHINE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATE_MACHINE, 0)
}

func (s *SourceContext) T_STATE_REPO() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATE_REPO, 0)
}

func (s *SourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SourceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSource(s)
	}
}

func (s *SourceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSource(s)
	}
}




func (p *SQLParser) Source() (localctx ISourceContext) {
	localctx = NewSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, SQLParserRULE_source)
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
		p.SetState(334)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserT_STATE_REPO || _la == SQLParserT_STATE_MACHINE) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
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
	p.EnterRule(localctx, 52, SQLParserRULE_queryStmt)
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
	p.SetState(337)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_EXPLAIN {
		{
			p.SetState(336)
			p.Match(SQLParserT_EXPLAIN)
		}

	}
	{
		p.SetState(339)
		p.SelectExpr()
	}
	{
		p.SetState(340)
		p.FromClause()
	}
	p.SetState(342)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(341)
			p.WhereClause()
		}

	}
	p.SetState(345)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_GROUP {
		{
			p.SetState(344)
			p.GroupByClause()
		}

	}
	p.SetState(348)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ORDER {
		{
			p.SetState(347)
			p.OrderByClause()
		}

	}
	p.SetState(351)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(350)
			p.LimitClause()
		}

	}
	p.SetState(354)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WITH_VALUE {
		{
			p.SetState(353)
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
	p.EnterRule(localctx, 54, SQLParserRULE_selectExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(356)
		p.Match(SQLParserT_SELECT)
	}
	{
		p.SetState(357)
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
	p.EnterRule(localctx, 56, SQLParserRULE_fields)
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
		p.SetState(359)
		p.Field()
	}
	p.SetState(364)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(360)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(361)
			p.Field()
		}


		p.SetState(366)
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
	p.EnterRule(localctx, 58, SQLParserRULE_field)
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
		p.SetState(367)
		p.fieldExpr(0)
	}
	p.SetState(369)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_AS {
		{
			p.SetState(368)
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
	p.EnterRule(localctx, 60, SQLParserRULE_alias)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(371)
		p.Match(SQLParserT_AS)
	}
	{
		p.SetState(372)
		p.Ident()
	}



	return localctx
}


// IStorageFilterContext is an interface to support dynamic dispatch.
type IStorageFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStorageFilterContext differentiates from other interfaces.
	IsStorageFilterContext()
}

type StorageFilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStorageFilterContext() *StorageFilterContext {
	var p = new(StorageFilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_storageFilter
	return p
}

func (*StorageFilterContext) IsStorageFilterContext() {}

func NewStorageFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StorageFilterContext {
	var p = new(StorageFilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_storageFilter

	return p
}

func (s *StorageFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *StorageFilterContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *StorageFilterContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *StorageFilterContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *StorageFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StorageFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StorageFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterStorageFilter(s)
	}
}

func (s *StorageFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitStorageFilter(s)
	}
}




func (p *SQLParser) StorageFilter() (localctx IStorageFilterContext) {
	localctx = NewStorageFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, SQLParserRULE_storageFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(374)
		p.Match(SQLParserT_STORAGE)
	}
	{
		p.SetState(375)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(376)
		p.Ident()
	}



	return localctx
}


// IDatabaseFilterContext is an interface to support dynamic dispatch.
type IDatabaseFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDatabaseFilterContext differentiates from other interfaces.
	IsDatabaseFilterContext()
}

type DatabaseFilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDatabaseFilterContext() *DatabaseFilterContext {
	var p = new(DatabaseFilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_databaseFilter
	return p
}

func (*DatabaseFilterContext) IsDatabaseFilterContext() {}

func NewDatabaseFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DatabaseFilterContext {
	var p = new(DatabaseFilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_databaseFilter

	return p
}

func (s *DatabaseFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *DatabaseFilterContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *DatabaseFilterContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *DatabaseFilterContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *DatabaseFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DatabaseFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *DatabaseFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDatabaseFilter(s)
	}
}

func (s *DatabaseFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDatabaseFilter(s)
	}
}




func (p *SQLParser) DatabaseFilter() (localctx IDatabaseFilterContext) {
	localctx = NewDatabaseFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, SQLParserRULE_databaseFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(378)
		p.Match(SQLParserT_DATASBAE)
	}
	{
		p.SetState(379)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(380)
		p.Ident()
	}



	return localctx
}


// ITypeFilterContext is an interface to support dynamic dispatch.
type ITypeFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeFilterContext differentiates from other interfaces.
	IsTypeFilterContext()
}

type TypeFilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeFilterContext() *TypeFilterContext {
	var p = new(TypeFilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_typeFilter
	return p
}

func (*TypeFilterContext) IsTypeFilterContext() {}

func NewTypeFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeFilterContext {
	var p = new(TypeFilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_typeFilter

	return p
}

func (s *TypeFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeFilterContext) T_TYPE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TYPE, 0)
}

func (s *TypeFilterContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *TypeFilterContext) Ident() IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *TypeFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TypeFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTypeFilter(s)
	}
}

func (s *TypeFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTypeFilter(s)
	}
}




func (p *SQLParser) TypeFilter() (localctx ITypeFilterContext) {
	localctx = NewTypeFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, SQLParserRULE_typeFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(382)
		p.Match(SQLParserT_TYPE)
	}
	{
		p.SetState(383)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(384)
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

func (s *FromClauseContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *FromClauseContext) Namespace() INamespaceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INamespaceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
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
	p.EnterRule(localctx, 68, SQLParserRULE_fromClause)
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
		p.SetState(386)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(387)
		p.MetricName()
	}
	p.SetState(390)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(388)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(389)
			p.Namespace()
		}

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
	p.EnterRule(localctx, 70, SQLParserRULE_whereClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(392)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(393)
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
	p.EnterRule(localctx, 72, SQLParserRULE_conditionExpr)
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

	p.SetState(405)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(395)
			p.tagFilterExpr(0)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(396)
			p.tagFilterExpr(0)
		}
		{
			p.SetState(397)
			p.Match(SQLParserT_AND)
		}
		{
			p.SetState(398)
			p.TimeRangeExpr()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(400)
			p.TimeRangeExpr()
		}
		p.SetState(403)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == SQLParserT_AND {
			{
				p.SetState(401)
				p.Match(SQLParserT_AND)
			}
			{
				p.SetState(402)
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
	_startState := 74
	p.EnterRecursionRule(localctx, 74, SQLParserRULE_tagFilterExpr, _p)
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
	p.SetState(435)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(408)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(409)
			p.tagFilterExpr(0)
		}
		{
			p.SetState(410)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(412)
			p.TagKey()
		}
		p.SetState(421)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_EQUAL:
			{
				p.SetState(413)
				p.Match(SQLParserT_EQUAL)
			}


		case SQLParserT_LIKE:
			{
				p.SetState(414)
				p.Match(SQLParserT_LIKE)
			}


		case SQLParserT_NOT:
			{
				p.SetState(415)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(416)
				p.Match(SQLParserT_LIKE)
			}


		case SQLParserT_REGEXP:
			{
				p.SetState(417)
				p.Match(SQLParserT_REGEXP)
			}


		case SQLParserT_NEQREGEXP:
			{
				p.SetState(418)
				p.Match(SQLParserT_NEQREGEXP)
			}


		case SQLParserT_NOTEQUAL:
			{
				p.SetState(419)
				p.Match(SQLParserT_NOTEQUAL)
			}


		case SQLParserT_NOTEQUAL2:
			{
				p.SetState(420)
				p.Match(SQLParserT_NOTEQUAL2)
			}



		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(423)
			p.TagValue()
		}


	case 3:
		{
			p.SetState(425)
			p.TagKey()
		}
		p.SetState(429)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_IN:
			{
				p.SetState(426)
				p.Match(SQLParserT_IN)
			}


		case SQLParserT_NOT:
			{
				p.SetState(427)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(428)
				p.Match(SQLParserT_IN)
			}



		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(431)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(432)
			p.TagValueList()
		}
		{
			p.SetState(433)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(442)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewTagFilterExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_tagFilterExpr)
			p.SetState(437)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(438)
				_la = p.GetTokenStream().LA(1)

				if !(_la == SQLParserT_AND || _la == SQLParserT_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(439)
				p.tagFilterExpr(2)
			}


		}
		p.SetState(444)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 76, SQLParserRULE_tagValueList)
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
		p.TagValue()
	}
	p.SetState(450)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(446)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(447)
			p.TagValue()
		}


		p.SetState(452)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// IMetricListFilterContext is an interface to support dynamic dispatch.
type IMetricListFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMetricListFilterContext differentiates from other interfaces.
	IsMetricListFilterContext()
}

type MetricListFilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetricListFilterContext() *MetricListFilterContext {
	var p = new(MetricListFilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_metricListFilter
	return p
}

func (*MetricListFilterContext) IsMetricListFilterContext() {}

func NewMetricListFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MetricListFilterContext {
	var p = new(MetricListFilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_metricListFilter

	return p
}

func (s *MetricListFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *MetricListFilterContext) T_METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRIC, 0)
}

func (s *MetricListFilterContext) T_IN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IN, 0)
}

func (s *MetricListFilterContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *MetricListFilterContext) MetricList() IMetricListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMetricListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMetricListContext)
}

func (s *MetricListFilterContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *MetricListFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MetricListFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *MetricListFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterMetricListFilter(s)
	}
}

func (s *MetricListFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitMetricListFilter(s)
	}
}




func (p *SQLParser) MetricListFilter() (localctx IMetricListFilterContext) {
	localctx = NewMetricListFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, SQLParserRULE_metricListFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
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
		p.Match(SQLParserT_METRIC)
	}
	{
		p.SetState(454)
		p.Match(SQLParserT_IN)
	}

	{
		p.SetState(455)
		p.Match(SQLParserT_OPEN_P)
	}
	{
		p.SetState(456)
		p.MetricList()
	}
	{
		p.SetState(457)
		p.Match(SQLParserT_CLOSE_P)
	}




	return localctx
}


// IMetricListContext is an interface to support dynamic dispatch.
type IMetricListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMetricListContext differentiates from other interfaces.
	IsMetricListContext()
}

type MetricListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetricListContext() *MetricListContext {
	var p = new(MetricListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_metricList
	return p
}

func (*MetricListContext) IsMetricListContext() {}

func NewMetricListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MetricListContext {
	var p = new(MetricListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_metricList

	return p
}

func (s *MetricListContext) GetParser() antlr.Parser { return s.parser }

func (s *MetricListContext) AllIdent() []IIdentContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IIdentContext)(nil)).Elem())
	var tst = make([]IIdentContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IIdentContext)
		}
	}

	return tst
}

func (s *MetricListContext) Ident(i int) IIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *MetricListContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *MetricListContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *MetricListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MetricListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *MetricListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterMetricList(s)
	}
}

func (s *MetricListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitMetricList(s)
	}
}




func (p *SQLParser) MetricList() (localctx IMetricListContext) {
	localctx = NewMetricListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, SQLParserRULE_metricList)
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
		p.SetState(459)
		p.Ident()
	}
	p.SetState(464)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(460)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(461)
			p.Ident()
		}


		p.SetState(466)
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
	p.EnterRule(localctx, 82, SQLParserRULE_timeRangeExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(467)
		p.TimeExpr()
	}
	p.SetState(470)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 31, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(468)
			p.Match(SQLParserT_AND)
		}
		{
			p.SetState(469)
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
	p.EnterRule(localctx, 84, SQLParserRULE_timeExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(472)
		p.Match(SQLParserT_TIME)
	}
	{
		p.SetState(473)
		p.BinaryOperator()
	}
	p.SetState(476)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 32, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(474)
			p.NowExpr()
		}


	case 2:
		{
			p.SetState(475)
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
	p.EnterRule(localctx, 86, SQLParserRULE_nowExpr)
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
		p.SetState(478)
		p.NowFunc()
	}
	p.SetState(480)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if ((((_la - 113)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 113))) & ((1 << (SQLParserT_ADD - 113)) | (1 << (SQLParserT_SUB - 113)) | (1 << (SQLParserL_INT - 113)))) != 0) {
		{
			p.SetState(479)
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
	p.EnterRule(localctx, 88, SQLParserRULE_nowFunc)
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
		p.SetState(482)
		p.Match(SQLParserT_NOW)
	}
	{
		p.SetState(483)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(485)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_USE) | (1 << SQLParserT_STATE_REPO) | (1 << SQLParserT_STATE_MACHINE) | (1 << SQLParserT_MASTER) | (1 << SQLParserT_METADATA) | (1 << SQLParserT_TYPES) | (1 << SQLParserT_TYPE) | (1 << SQLParserT_STORAGES) | (1 << SQLParserT_STORAGE) | (1 << SQLParserT_BROKER) | (1 << SQLParserT_ALIVE))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_SCHEMAS - 32)) | (1 << (SQLParserT_DATASBAE - 32)) | (1 << (SQLParserT_DATASBAES - 32)) | (1 << (SQLParserT_NAMESPACE - 32)) | (1 << (SQLParserT_NAMESPACES - 32)) | (1 << (SQLParserT_NODE - 32)) | (1 << (SQLParserT_METRICS - 32)) | (1 << (SQLParserT_METRIC - 32)) | (1 << (SQLParserT_FIELD - 32)) | (1 << (SQLParserT_FIELDS - 32)) | (1 << (SQLParserT_TAG - 32)) | (1 << (SQLParserT_INFO - 32)) | (1 << (SQLParserT_KEYS - 32)) | (1 << (SQLParserT_KEY - 32)) | (1 << (SQLParserT_WITH - 32)) | (1 << (SQLParserT_VALUES - 32)) | (1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_ASC - 64)) | (1 << (SQLParserT_DESC - 64)) | (1 << (SQLParserT_LIKE - 64)) | (1 << (SQLParserT_NOT - 64)) | (1 << (SQLParserT_BETWEEN - 64)) | (1 << (SQLParserT_IS - 64)) | (1 << (SQLParserT_GROUP - 64)) | (1 << (SQLParserT_HAVING - 64)) | (1 << (SQLParserT_BY - 64)) | (1 << (SQLParserT_FOR - 64)) | (1 << (SQLParserT_STATS - 64)) | (1 << (SQLParserT_TIME - 64)) | (1 << (SQLParserT_NOW - 64)) | (1 << (SQLParserT_IN - 64)) | (1 << (SQLParserT_LOG - 64)) | (1 << (SQLParserT_PROFILE - 64)) | (1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_RATE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)))) != 0) || ((((_la - 111)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 111))) & ((1 << (SQLParserT_OPEN_P - 111)) | (1 << (SQLParserT_ADD - 111)) | (1 << (SQLParserT_SUB - 111)) | (1 << (SQLParserL_ID - 111)) | (1 << (SQLParserL_INT - 111)) | (1 << (SQLParserL_DEC - 111)))) != 0) {
		{
			p.SetState(484)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(487)
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
	p.EnterRule(localctx, 90, SQLParserRULE_groupByClause)
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
		p.SetState(489)
		p.Match(SQLParserT_GROUP)
	}
	{
		p.SetState(490)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(491)
		p.GroupByKeys()
	}
	p.SetState(497)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_FILL {
		{
			p.SetState(492)
			p.Match(SQLParserT_FILL)
		}
		{
			p.SetState(493)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(494)
			p.FillOption()
		}
		{
			p.SetState(495)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.SetState(500)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_HAVING {
		{
			p.SetState(499)
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
	p.EnterRule(localctx, 92, SQLParserRULE_groupByKeys)
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
		p.GroupByKey()
	}
	p.SetState(507)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(503)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(504)
			p.GroupByKey()
		}


		p.SetState(509)
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
	p.EnterRule(localctx, 94, SQLParserRULE_groupByKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(516)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(510)
			p.Ident()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(511)
			p.Match(SQLParserT_TIME)
		}
		{
			p.SetState(512)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(513)
			p.DurationLit()
		}
		{
			p.SetState(514)
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
	p.EnterRule(localctx, 96, SQLParserRULE_fillOption)
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
		p.SetState(518)
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
	p.EnterRule(localctx, 98, SQLParserRULE_orderByClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(520)
		p.Match(SQLParserT_ORDER)
	}
	{
		p.SetState(521)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(522)
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
	p.EnterRule(localctx, 100, SQLParserRULE_sortField)
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
		p.SetState(524)
		p.fieldExpr(0)
	}
	p.SetState(528)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_ASC || _la == SQLParserT_DESC {
		{
			p.SetState(525)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_ASC || _la == SQLParserT_DESC) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}


		p.SetState(530)
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
	p.EnterRule(localctx, 102, SQLParserRULE_sortFields)
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
		p.SetState(531)
		p.SortField()
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
			p.SortField()
		}


		p.SetState(538)
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
	p.EnterRule(localctx, 104, SQLParserRULE_havingClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(539)
		p.Match(SQLParserT_HAVING)
	}
	{
		p.SetState(540)
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
	_startState := 106
	p.EnterRecursionRule(localctx, 106, SQLParserRULE_boolExpr, _p)

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
	p.SetState(548)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(543)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(544)
			p.boolExpr(0)
		}
		{
			p.SetState(545)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(547)
			p.BoolExprAtom()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(556)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewBoolExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_boolExpr)
			p.SetState(550)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(551)
				p.BoolExprLogicalOp()
			}
			{
				p.SetState(552)
				p.boolExpr(3)
			}


		}
		p.SetState(558)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 108, SQLParserRULE_boolExprLogicalOp)
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
		p.SetState(559)
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
	p.EnterRule(localctx, 110, SQLParserRULE_boolExprAtom)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(561)
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
	p.EnterRule(localctx, 112, SQLParserRULE_binaryExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(563)
		p.fieldExpr(0)
	}
	{
		p.SetState(564)
		p.BinaryOperator()
	}
	{
		p.SetState(565)
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
	p.EnterRule(localctx, 114, SQLParserRULE_binaryOperator)
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

	p.SetState(575)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_EQUAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(567)
			p.Match(SQLParserT_EQUAL)
		}


	case SQLParserT_NOTEQUAL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(568)
			p.Match(SQLParserT_NOTEQUAL)
		}


	case SQLParserT_NOTEQUAL2:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(569)
			p.Match(SQLParserT_NOTEQUAL2)
		}


	case SQLParserT_LESS:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(570)
			p.Match(SQLParserT_LESS)
		}


	case SQLParserT_LESSEQUAL:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(571)
			p.Match(SQLParserT_LESSEQUAL)
		}


	case SQLParserT_GREATER:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(572)
			p.Match(SQLParserT_GREATER)
		}


	case SQLParserT_GREATEREQUAL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(573)
			p.Match(SQLParserT_GREATEREQUAL)
		}


	case SQLParserT_LIKE, SQLParserT_REGEXP:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(574)
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
	_startState := 116
	p.EnterRecursionRule(localctx, 116, SQLParserRULE_fieldExpr, _p)

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
	p.SetState(585)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 44, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(578)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(579)
			p.fieldExpr(0)
		}
		{
			p.SetState(580)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(582)
			p.ExprFunc()
		}


	case 3:
		{
			p.SetState(583)
			p.ExprAtom()
		}


	case 4:
		{
			p.SetState(584)
			p.DurationLit()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(601)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(599)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext()) {
			case 1:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(587)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(588)
					p.Match(SQLParserT_MUL)
				}
				{
					p.SetState(589)
					p.fieldExpr(9)
				}


			case 2:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(590)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(591)
					p.Match(SQLParserT_DIV)
				}
				{
					p.SetState(592)
					p.fieldExpr(8)
				}


			case 3:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(593)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				}
				{
					p.SetState(594)
					p.Match(SQLParserT_ADD)
				}
				{
					p.SetState(595)
					p.fieldExpr(7)
				}


			case 4:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(596)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				{
					p.SetState(597)
					p.Match(SQLParserT_SUB)
				}
				{
					p.SetState(598)
					p.fieldExpr(6)
				}

			}

		}
		p.SetState(603)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 118, SQLParserRULE_durationLit)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(604)
		p.IntNumber()
	}
	{
		p.SetState(605)
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
	p.EnterRule(localctx, 120, SQLParserRULE_intervalItem)
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

		if !(((((_la - 88)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 88))) & ((1 << (SQLParserT_SECOND - 88)) | (1 << (SQLParserT_MINUTE - 88)) | (1 << (SQLParserT_HOUR - 88)) | (1 << (SQLParserT_DAY - 88)) | (1 << (SQLParserT_WEEK - 88)) | (1 << (SQLParserT_MONTH - 88)) | (1 << (SQLParserT_YEAR - 88)))) != 0)) {
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
	p.EnterRule(localctx, 122, SQLParserRULE_exprFunc)
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
		p.SetState(609)
		p.FuncName()
	}
	{
		p.SetState(610)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(612)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_USE) | (1 << SQLParserT_STATE_REPO) | (1 << SQLParserT_STATE_MACHINE) | (1 << SQLParserT_MASTER) | (1 << SQLParserT_METADATA) | (1 << SQLParserT_TYPES) | (1 << SQLParserT_TYPE) | (1 << SQLParserT_STORAGES) | (1 << SQLParserT_STORAGE) | (1 << SQLParserT_BROKER) | (1 << SQLParserT_ALIVE))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_SCHEMAS - 32)) | (1 << (SQLParserT_DATASBAE - 32)) | (1 << (SQLParserT_DATASBAES - 32)) | (1 << (SQLParserT_NAMESPACE - 32)) | (1 << (SQLParserT_NAMESPACES - 32)) | (1 << (SQLParserT_NODE - 32)) | (1 << (SQLParserT_METRICS - 32)) | (1 << (SQLParserT_METRIC - 32)) | (1 << (SQLParserT_FIELD - 32)) | (1 << (SQLParserT_FIELDS - 32)) | (1 << (SQLParserT_TAG - 32)) | (1 << (SQLParserT_INFO - 32)) | (1 << (SQLParserT_KEYS - 32)) | (1 << (SQLParserT_KEY - 32)) | (1 << (SQLParserT_WITH - 32)) | (1 << (SQLParserT_VALUES - 32)) | (1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_ASC - 64)) | (1 << (SQLParserT_DESC - 64)) | (1 << (SQLParserT_LIKE - 64)) | (1 << (SQLParserT_NOT - 64)) | (1 << (SQLParserT_BETWEEN - 64)) | (1 << (SQLParserT_IS - 64)) | (1 << (SQLParserT_GROUP - 64)) | (1 << (SQLParserT_HAVING - 64)) | (1 << (SQLParserT_BY - 64)) | (1 << (SQLParserT_FOR - 64)) | (1 << (SQLParserT_STATS - 64)) | (1 << (SQLParserT_TIME - 64)) | (1 << (SQLParserT_NOW - 64)) | (1 << (SQLParserT_IN - 64)) | (1 << (SQLParserT_LOG - 64)) | (1 << (SQLParserT_PROFILE - 64)) | (1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_RATE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)))) != 0) || ((((_la - 111)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 111))) & ((1 << (SQLParserT_OPEN_P - 111)) | (1 << (SQLParserT_ADD - 111)) | (1 << (SQLParserT_SUB - 111)) | (1 << (SQLParserL_ID - 111)) | (1 << (SQLParserL_INT - 111)) | (1 << (SQLParserL_DEC - 111)))) != 0) {
		{
			p.SetState(611)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(614)
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

func (s *FuncNameContext) T_RATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_RATE, 0)
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
	p.EnterRule(localctx, 124, SQLParserRULE_funcName)
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
		p.SetState(616)
		_la = p.GetTokenStream().LA(1)

		if !(((((_la - 80)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 80))) & ((1 << (SQLParserT_SUM - 80)) | (1 << (SQLParserT_MIN - 80)) | (1 << (SQLParserT_MAX - 80)) | (1 << (SQLParserT_COUNT - 80)) | (1 << (SQLParserT_AVG - 80)) | (1 << (SQLParserT_STDDEV - 80)) | (1 << (SQLParserT_QUANTILE - 80)) | (1 << (SQLParserT_RATE - 80)))) != 0)) {
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
	p.EnterRule(localctx, 126, SQLParserRULE_exprFuncParams)
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
		p.SetState(618)
		p.FuncParam()
	}
	p.SetState(623)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(619)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(620)
			p.FuncParam()
		}


		p.SetState(625)
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
	p.EnterRule(localctx, 128, SQLParserRULE_funcParam)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(628)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 49, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(626)
			p.fieldExpr(0)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(627)
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
	p.EnterRule(localctx, 130, SQLParserRULE_exprAtom)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(636)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(630)
			p.Ident()
		}
		p.SetState(632)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 50, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(631)
				p.IdentFilter()
			}


		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(634)
			p.DecNumber()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(635)
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
	p.EnterRule(localctx, 132, SQLParserRULE_identFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(638)
		p.Match(SQLParserT_OPEN_SB)
	}
	{
		p.SetState(639)
		p.tagFilterExpr(0)
	}
	{
		p.SetState(640)
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
	p.EnterRule(localctx, 134, SQLParserRULE_json)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(642)
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
	p.EnterRule(localctx, 136, SQLParserRULE_obj)
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

	p.SetState(657)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 53, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(644)
			p.Match(SQLParserT_OPEN_B)
		}
		{
			p.SetState(645)
			p.Pair()
		}
		p.SetState(650)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for _la == SQLParserT_COMMA {
			{
				p.SetState(646)
				p.Match(SQLParserT_COMMA)
			}
			{
				p.SetState(647)
				p.Pair()
			}


			p.SetState(652)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(653)
			p.Match(SQLParserT_CLOSE_B)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(655)
			p.Match(SQLParserT_OPEN_B)
		}
		{
			p.SetState(656)
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
	p.EnterRule(localctx, 138, SQLParserRULE_pair)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(659)
		p.Match(SQLParserSTRING)
	}
	{
		p.SetState(660)
		p.Match(SQLParserT_COLON)
	}
	{
		p.SetState(661)
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
	p.EnterRule(localctx, 140, SQLParserRULE_arr)
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

	p.SetState(676)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 55, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(663)
			p.Match(SQLParserT_OPEN_SB)
		}
		{
			p.SetState(664)
			p.Value()
		}
		p.SetState(669)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for _la == SQLParserT_COMMA {
			{
				p.SetState(665)
				p.Match(SQLParserT_COMMA)
			}
			{
				p.SetState(666)
				p.Value()
			}


			p.SetState(671)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(672)
			p.Match(SQLParserT_CLOSE_SB)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(674)
			p.Match(SQLParserT_OPEN_SB)
		}
		{
			p.SetState(675)
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
	p.EnterRule(localctx, 142, SQLParserRULE_value)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(686)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 56, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(678)
			p.Match(SQLParserSTRING)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(679)
			p.IntNumber()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(680)
			p.DecNumber()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(681)
			p.Obj()
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(682)
			p.Arr()
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(683)
			p.Match(SQLParserT__0)
		}


	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(684)
			p.Match(SQLParserT__1)
		}


	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(685)
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
	p.SetState(689)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(688)
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
		p.SetState(691)
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
	p.SetState(694)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(693)
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
		p.SetState(696)
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
		p.SetState(698)
		p.Match(SQLParserT_LIMIT)
	}
	{
		p.SetState(699)
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
		p.SetState(701)
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
		p.SetState(703)
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
		p.SetState(705)
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
	p.EnterRule(localctx, 156, SQLParserRULE_ident)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
	p.SetState(709)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserL_ID:
		{
			p.SetState(707)
			p.Match(SQLParserL_ID)
		}


	case SQLParserT_CREATE, SQLParserT_UPDATE, SQLParserT_SET, SQLParserT_DROP, SQLParserT_INTERVAL, SQLParserT_INTERVAL_NAME, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_META_TTL, SQLParserT_PAST_TTL, SQLParserT_FUTURE_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_USE, SQLParserT_STATE_REPO, SQLParserT_STATE_MACHINE, SQLParserT_MASTER, SQLParserT_METADATA, SQLParserT_TYPES, SQLParserT_TYPE, SQLParserT_STORAGES, SQLParserT_STORAGE, SQLParserT_BROKER, SQLParserT_ALIVE, SQLParserT_SCHEMAS, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NAMESPACE, SQLParserT_NAMESPACES, SQLParserT_NODE, SQLParserT_METRICS, SQLParserT_METRIC, SQLParserT_FIELD, SQLParserT_FIELDS, SQLParserT_TAG, SQLParserT_INFO, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_VALUE, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_EXPLAIN, SQLParserT_WITH_VALUE, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_HAVING, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_NOW, SQLParserT_IN, SQLParserT_LOG, SQLParserT_PROFILE, SQLParserT_SUM, SQLParserT_MIN, SQLParserT_MAX, SQLParserT_COUNT, SQLParserT_AVG, SQLParserT_STDDEV, SQLParserT_QUANTILE, SQLParserT_RATE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
		{
			p.SetState(708)
			p.NonReservedWords()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(718)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 61, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(711)
				p.Match(SQLParserT_DOT)
			}
			p.SetState(714)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case SQLParserL_ID:
				{
					p.SetState(712)
					p.Match(SQLParserL_ID)
				}


			case SQLParserT_CREATE, SQLParserT_UPDATE, SQLParserT_SET, SQLParserT_DROP, SQLParserT_INTERVAL, SQLParserT_INTERVAL_NAME, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_META_TTL, SQLParserT_PAST_TTL, SQLParserT_FUTURE_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_USE, SQLParserT_STATE_REPO, SQLParserT_STATE_MACHINE, SQLParserT_MASTER, SQLParserT_METADATA, SQLParserT_TYPES, SQLParserT_TYPE, SQLParserT_STORAGES, SQLParserT_STORAGE, SQLParserT_BROKER, SQLParserT_ALIVE, SQLParserT_SCHEMAS, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NAMESPACE, SQLParserT_NAMESPACES, SQLParserT_NODE, SQLParserT_METRICS, SQLParserT_METRIC, SQLParserT_FIELD, SQLParserT_FIELDS, SQLParserT_TAG, SQLParserT_INFO, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_VALUE, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_EXPLAIN, SQLParserT_WITH_VALUE, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_HAVING, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_NOW, SQLParserT_IN, SQLParserT_LOG, SQLParserT_PROFILE, SQLParserT_SUM, SQLParserT_MIN, SQLParserT_MAX, SQLParserT_COUNT, SQLParserT_AVG, SQLParserT_STDDEV, SQLParserT_QUANTILE, SQLParserT_RATE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
				{
					p.SetState(713)
					p.NonReservedWords()
				}



			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}


		}
		p.SetState(720)
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

func (s *NonReservedWordsContext) T_RATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_RATE, 0)
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

func (s *NonReservedWordsContext) T_METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METADATA, 0)
}

func (s *NonReservedWordsContext) T_TYPE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TYPE, 0)
}

func (s *NonReservedWordsContext) T_TYPES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TYPES, 0)
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

func (s *NonReservedWordsContext) T_STATE_REPO() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATE_REPO, 0)
}

func (s *NonReservedWordsContext) T_STATE_MACHINE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATE_MACHINE, 0)
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
	p.EnterRule(localctx, 158, SQLParserRULE_nonReservedWords)
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
		p.SetState(721)
		_la = p.GetTokenStream().LA(1)

		if !((((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_USE) | (1 << SQLParserT_STATE_REPO) | (1 << SQLParserT_STATE_MACHINE) | (1 << SQLParserT_MASTER) | (1 << SQLParserT_METADATA) | (1 << SQLParserT_TYPES) | (1 << SQLParserT_TYPE) | (1 << SQLParserT_STORAGES) | (1 << SQLParserT_STORAGE) | (1 << SQLParserT_BROKER) | (1 << SQLParserT_ALIVE))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_SCHEMAS - 32)) | (1 << (SQLParserT_DATASBAE - 32)) | (1 << (SQLParserT_DATASBAES - 32)) | (1 << (SQLParserT_NAMESPACE - 32)) | (1 << (SQLParserT_NAMESPACES - 32)) | (1 << (SQLParserT_NODE - 32)) | (1 << (SQLParserT_METRICS - 32)) | (1 << (SQLParserT_METRIC - 32)) | (1 << (SQLParserT_FIELD - 32)) | (1 << (SQLParserT_FIELDS - 32)) | (1 << (SQLParserT_TAG - 32)) | (1 << (SQLParserT_INFO - 32)) | (1 << (SQLParserT_KEYS - 32)) | (1 << (SQLParserT_KEY - 32)) | (1 << (SQLParserT_WITH - 32)) | (1 << (SQLParserT_VALUES - 32)) | (1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_ASC - 64)) | (1 << (SQLParserT_DESC - 64)) | (1 << (SQLParserT_LIKE - 64)) | (1 << (SQLParserT_NOT - 64)) | (1 << (SQLParserT_BETWEEN - 64)) | (1 << (SQLParserT_IS - 64)) | (1 << (SQLParserT_GROUP - 64)) | (1 << (SQLParserT_HAVING - 64)) | (1 << (SQLParserT_BY - 64)) | (1 << (SQLParserT_FOR - 64)) | (1 << (SQLParserT_STATS - 64)) | (1 << (SQLParserT_TIME - 64)) | (1 << (SQLParserT_NOW - 64)) | (1 << (SQLParserT_IN - 64)) | (1 << (SQLParserT_LOG - 64)) | (1 << (SQLParserT_PROFILE - 64)) | (1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_RATE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)) | (1 << (SQLParserT_YEAR - 64)))) != 0)) {
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
	case 37:
			var t *TagFilterExprContext = nil
			if localctx != nil { t = localctx.(*TagFilterExprContext) }
			return p.TagFilterExpr_Sempred(t, predIndex)

	case 53:
			var t *BoolExprContext = nil
			if localctx != nil { t = localctx.(*BoolExprContext) }
			return p.BoolExpr_Sempred(t, predIndex)

	case 58:
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

