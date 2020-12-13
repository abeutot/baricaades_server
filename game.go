package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

/* game states:
0 WAITING FOR PLAYERS
1 PLAYING
2 FINISHED
*/

/* TODO concurrency */
type Game struct {
	players   []string
	state     uint8
	dice      uint8
	positions [31]uint8
	/* positions are RED (5), GREEN(5), YELLOW (5), BLUE(5), BARICADES(11) in that order */
	subscriptions []chan *GameJson
}

const (
	POSITION_I16 = 0
	POSITION_A15 = 1
	POSITION_B15 = 2
	POSITION_C15 = 3
	POSITION_D15 = 4
	POSITION_E15 = 5
	POSITION_F15 = 6
	POSITION_G15 = 7
	POSITION_H15 = 8
	POSITION_I15 = 9
	POSITION_J15 = 10
	POSITION_K15 = 11
	POSITION_L15 = 12
	POSITION_M15 = 13
	POSITION_N15 = 14
	POSITION_O15 = 15
	POSITION_P15 = 16
	POSITION_Q15 = 17
	POSITION_A14 = 18
	POSITION_Q14 = 19
	POSITION_A13 = 20
	POSITION_B13 = 21
	POSITION_C13 = 22
	POSITION_D13 = 23
	POSITION_E13 = 24
	POSITION_F13 = 25
	POSITION_G13 = 26
	POSITION_H13 = 27
	POSITION_I13 = 28
	POSITION_J13 = 29
	POSITION_K13 = 30
	POSITION_L13 = 31
	POSITION_M13 = 32
	POSITION_N13 = 33
	POSITION_O13 = 34
	POSITION_P13 = 35
	POSITION_Q13 = 36
	POSITION_I12 = 37
	POSITION_G11 = 38
	POSITION_H11 = 39
	POSITION_I11 = 40
	POSITION_J11 = 41
	POSITION_K11 = 42
	POSITION_G10 = 43
	POSITION_K10 = 44
	POSITION_E9  = 45
	POSITION_F9  = 46
	POSITION_G9  = 47
	POSITION_H9  = 48
	POSITION_I9  = 49
	POSITION_J9  = 50
	POSITION_K9  = 51
	POSITION_L9  = 52
	POSITION_M9  = 53
	POSITION_E8  = 54
	POSITION_M8  = 55
	POSITION_C7  = 56
	POSITION_D7  = 57
	POSITION_E7  = 58
	POSITION_F7  = 59
	POSITION_G7  = 60
	POSITION_H7  = 61
	POSITION_I7  = 62
	POSITION_J7  = 63
	POSITION_K7  = 64
	POSITION_L7  = 65
	POSITION_M7  = 66
	POSITION_N7  = 67
	POSITION_O7  = 68
	POSITION_C6  = 69
	POSITION_O6  = 70
	POSITION_A5  = 71
	POSITION_B5  = 72
	POSITION_C5  = 73
	POSITION_D5  = 74
	POSITION_E5  = 75
	POSITION_F5  = 76
	POSITION_G5  = 77
	POSITION_H5  = 78
	POSITION_I5  = 79
	POSITION_J5  = 80
	POSITION_K5  = 81
	POSITION_L5  = 82
	POSITION_M5  = 83
	POSITION_N5  = 84
	POSITION_O5  = 85
	POSITION_P5  = 86
	POSITION_Q5  = 87
	POSITION_A4  = 88
	POSITION_E4  = 89
	POSITION_I4  = 90
	POSITION_M4  = 91
	POSITION_Q4  = 92
	POSITION_A3  = 93
	POSITION_B3  = 94
	POSITION_C3  = 95
	POSITION_D3  = 96
	POSITION_E3  = 97
	POSITION_F3  = 98
	POSITION_G3  = 99
	POSITION_H3  = 100
	POSITION_I3  = 101
	POSITION_J3  = 102
	POSITION_K3  = 103
	POSITION_L3  = 104
	POSITION_M3  = 105
	POSITION_N3  = 106
	POSITION_O3  = 107
	POSITION_P3  = 108
	POSITION_Q3  = 109
	POSITION_B2  = 110
	POSITION_C2  = 111
	POSITION_D2  = 112
	POSITION_F2  = 113
	POSITION_G2  = 114
	POSITION_H2  = 115
	POSITION_J2  = 116
	POSITION_K2  = 117
	POSITION_L2  = 118
	POSITION_N2  = 119
	POSITION_O2  = 120
	POSITION_P2  = 121
	POSITION_B1  = 122
	POSITION_D1  = 123
	POSITION_F1  = 124
	POSITION_H1  = 125
	POSITION_J1  = 126
	POSITION_L1  = 127
	POSITION_N1  = 128
	POSITION_P1  = 129
)

var indexToPosition = []string{
	"I16",
	"A15",
	"B15",
	"C15",
	"D15",
	"E15",
	"F15",
	"G15",
	"H15",
	"I15",
	"J15",
	"K15",
	"L15",
	"M15",
	"N15",
	"O15",
	"P15",
	"Q15",
	"A14",
	"Q14",
	"A13",
	"B13",
	"C13",
	"D13",
	"E13",
	"F13",
	"G13",
	"H13",
	"I13",
	"J13",
	"K13",
	"L13",
	"M13",
	"N13",
	"O13",
	"P13",
	"Q13",
	"I12",
	"G11",
	"H11",
	"I11",
	"J11",
	"K11",
	"G10",
	"K10",
	"E9",
	"F9",
	"G9",
	"H9",
	"I9",
	"J9",
	"K9",
	"L9",
	"M9",
	"E8",
	"M8",
	"C7",
	"D7",
	"E7",
	"F7",
	"G7",
	"H7",
	"I7",
	"J7",
	"K7",
	"L7",
	"M7",
	"N7",
	"O7",
	"C6",
	"O6",
	"A5",
	"B5",
	"C5",
	"D5",
	"E5",
	"F5",
	"G5",
	"H5",
	"I5",
	"J5",
	"K5",
	"L5",
	"M5",
	"N5",
	"O5",
	"P5",
	"Q5",
	"A4",
	"E4",
	"I4",
	"M4",
	"Q4",
	"A3",
	"B3",
	"C3",
	"D3",
	"E3",
	"F3",
	"G3",
	"H3",
	"I3",
	"J3",
	"K3",
	"L3",
	"M3",
	"N3",
	"O3",
	"P3",
	"Q3",
	"B2",
	"C2",
	"D2",
	"F2",
	"G2",
	"H2",
	"J2",
	"K2",
	"L2",
	"N2",
	"O2",
	"P2",
	"B1",
	"D1",
	"F1",
	"H1",
	"J1",
	"L1",
	"N1",
	"P1",
}

var positionToIndex = map[string]uint8{
	"I16": 0,
	"A15": 1,
	"B15": 2,
	"C15": 3,
	"D15": 4,
	"E15": 5,
	"F15": 6,
	"G15": 7,
	"H15": 8,
	"I15": 9,
	"J15": 10,
	"K15": 11,
	"L15": 12,
	"M15": 13,
	"N15": 14,
	"O15": 15,
	"P15": 16,
	"Q15": 17,
	"A14": 18,
	"Q14": 19,
	"A13": 20,
	"B13": 21,
	"C13": 22,
	"D13": 23,
	"E13": 24,
	"F13": 25,
	"G13": 26,
	"H13": 27,
	"I13": 28,
	"J13": 29,
	"K13": 30,
	"L13": 31,
	"M13": 32,
	"N13": 33,
	"O13": 34,
	"P13": 35,
	"Q13": 36,
	"I12": 37,
	"G11": 38,
	"H11": 39,
	"I11": 40,
	"J11": 41,
	"K11": 42,
	"G10": 43,
	"K10": 44,
	"E9":  45,
	"F9":  46,
	"G9":  47,
	"H9":  48,
	"I9":  49,
	"J9":  50,
	"K9":  51,
	"L9":  52,
	"M9":  53,
	"E8":  54,
	"M8":  55,
	"C7":  56,
	"D7":  57,
	"E7":  58,
	"F7":  59,
	"G7":  60,
	"H7":  61,
	"I7":  62,
	"J7":  63,
	"K7":  64,
	"L7":  65,
	"M7":  66,
	"N7":  67,
	"O7":  68,
	"C6":  69,
	"O6":  70,
	"A5":  71,
	"B5":  72,
	"C5":  73,
	"D5":  74,
	"E5":  75,
	"F5":  76,
	"G5":  77,
	"H5":  78,
	"I5":  79,
	"J5":  80,
	"K5":  81,
	"L5":  82,
	"M5":  83,
	"N5":  84,
	"O5":  85,
	"P5":  86,
	"Q5":  87,
	"A4":  88,
	"E4":  89,
	"I4":  90,
	"M4":  91,
	"Q4":  92,
	"A3":  93,
	"B3":  94,
	"C3":  95,
	"D3":  96,
	"E3":  97,
	"F3":  98,
	"G3":  99,
	"H3":  100,
	"I3":  101,
	"J3":  102,
	"K3":  103,
	"L3":  104,
	"M3":  105,
	"N3":  106,
	"O3":  107,
	"P3":  108,
	"Q3":  109,
	"B2":  110,
	"C2":  111,
	"D2":  112,
	"F2":  113,
	"G2":  114,
	"H2":  115,
	"J2":  116,
	"K2":  117,
	"L2":  118,
	"N2":  119,
	"O2":  120,
	"P2":  121,
	"B1":  122,
	"D1":  123,
	"F1":  124,
	"H1":  125,
	"J1":  126,
	"L1":  127,
	"N1":  128,
	"P1":  129,
}

const (
	/* pawns types */
	PAWN_RED      uint8 = 1
	PAWN_GREEN    uint8 = 2
	PAWN_YELLOW   uint8 = 4
	PAWN_BLUE     uint8 = 8
	PAWN_BARICADE uint8 = 16
	PAWN_PLAYER   uint8 = PAWN_RED | PAWN_GREEN | PAWN_YELLOW | PAWN_BLUE
	PAWN_ALL      uint8 = PAWN_PLAYER | PAWN_BARICADE
)

var allowedPawns = [...]uint8{
	PAWN_PLAYER, // I16
	PAWN_ALL,    // A15
	PAWN_ALL,    // B15
	PAWN_ALL,    // C15
	PAWN_ALL,    // D15
	PAWN_ALL,    // E15
	PAWN_ALL,    // F15
	PAWN_ALL,    // G15
	PAWN_ALL,    // H15
	PAWN_ALL,    // I15
	PAWN_ALL,    // J15
	PAWN_ALL,    // K15
	PAWN_ALL,    // L15
	PAWN_ALL,    // M15
	PAWN_ALL,    // N15
	PAWN_ALL,    // O15
	PAWN_ALL,    // P15
	PAWN_ALL,    // Q15
	PAWN_ALL,    // A14
	PAWN_ALL,    // Q14
	PAWN_ALL,    // A13
	PAWN_ALL,    // B13
	PAWN_ALL,    // C13
	PAWN_ALL,    // D13
	PAWN_ALL,    // E13
	PAWN_ALL,    // F13
	PAWN_ALL,    // G13
	PAWN_ALL,    // H13
	PAWN_ALL,    // I13
	PAWN_ALL,    // J13
	PAWN_ALL,    // K13
	PAWN_ALL,    // L13
	PAWN_ALL,    // M13
	PAWN_ALL,    // N13
	PAWN_ALL,    // O13
	PAWN_ALL,    // P13
	PAWN_ALL,    // Q13
	PAWN_ALL,    // I12
	PAWN_ALL,    // G11
	PAWN_ALL,    // H11
	PAWN_ALL,    // I11
	PAWN_ALL,    // J11
	PAWN_ALL,    // K11
	PAWN_ALL,    // G10
	PAWN_ALL,    // K10
	PAWN_ALL,    // E9
	PAWN_ALL,    // F9
	PAWN_ALL,    // G9
	PAWN_ALL,    // H9
	PAWN_ALL,    // I9
	PAWN_ALL,    // J9
	PAWN_ALL,    // K9
	PAWN_ALL,    // L9
	PAWN_ALL,    // M9
	PAWN_ALL,    // E8
	PAWN_ALL,    // M8
	PAWN_ALL,    // C7
	PAWN_ALL,    // D7
	PAWN_ALL,    // E7
	PAWN_ALL,    // F7
	PAWN_ALL,    // G7
	PAWN_ALL,    // H7
	PAWN_ALL,    // I7
	PAWN_ALL,    // J7
	PAWN_ALL,    // K7
	PAWN_ALL,    // L7
	PAWN_ALL,    // M7
	PAWN_ALL,    // N7
	PAWN_ALL,    // O7
	PAWN_ALL,    // C6
	PAWN_ALL,    // O6
	PAWN_ALL,    // A5
	PAWN_ALL,    // B5
	PAWN_ALL,    // C5
	PAWN_ALL,    // D5
	PAWN_ALL,    // E5
	PAWN_ALL,    // F5
	PAWN_ALL,    // G5
	PAWN_ALL,    // H5
	PAWN_ALL,    // I5
	PAWN_ALL,    // J5
	PAWN_ALL,    // K5
	PAWN_ALL,    // L5
	PAWN_ALL,    // M5
	PAWN_ALL,    // N5
	PAWN_ALL,    // O5
	PAWN_ALL,    // P5
	PAWN_ALL,    // Q5
	PAWN_ALL,    // A4
	PAWN_ALL,    // E4
	PAWN_ALL,    // I4
	PAWN_ALL,    // M4
	PAWN_ALL,    // Q4
	PAWN_PLAYER, // A3
	PAWN_PLAYER, // B3
	PAWN_PLAYER, // C3
	PAWN_PLAYER, // D3
	PAWN_PLAYER, // E3
	PAWN_PLAYER, // F3
	PAWN_PLAYER, // G3
	PAWN_PLAYER, // H3
	PAWN_PLAYER, // I3
	PAWN_PLAYER, // J3
	PAWN_PLAYER, // K3
	PAWN_PLAYER, // L3
	PAWN_PLAYER, // M3
	PAWN_PLAYER, // N3
	PAWN_PLAYER, // O3
	PAWN_PLAYER, // P3
	PAWN_PLAYER, // Q3
	PAWN_RED,    // B2
	PAWN_RED,    // C2
	PAWN_RED,    // D2
	PAWN_GREEN,  // F2
	PAWN_GREEN,  // G2
	PAWN_GREEN,  // H2
	PAWN_YELLOW, // J2
	PAWN_YELLOW, // K2
	PAWN_YELLOW, // L2
	PAWN_BLUE,   // N2
	PAWN_BLUE,   // O2
	PAWN_BLUE,   // P2
	PAWN_RED,    // B1
	PAWN_RED,    // D1
	PAWN_GREEN,  // F1
	PAWN_GREEN,  // H1
	PAWN_YELLOW, // J1
	PAWN_YELLOW, // L1
	PAWN_BLUE,   // N1
	PAWN_BLUE,   // P1
}

/* those represent one way transitions from one position to another, in other
* words this is the graph representing the board */
var transitions = func() [130][130]bool {
	ret := [130][130]bool{}

	/* start transitions */
	ret[POSITION_B1][POSITION_C3] = true
	ret[POSITION_D1][POSITION_C3] = true
	ret[POSITION_B2][POSITION_C3] = true
	ret[POSITION_C2][POSITION_C3] = true
	ret[POSITION_D2][POSITION_C3] = true

	ret[POSITION_F1][POSITION_G3] = true
	ret[POSITION_H1][POSITION_G3] = true
	ret[POSITION_F2][POSITION_G3] = true
	ret[POSITION_G2][POSITION_G3] = true
	ret[POSITION_H2][POSITION_G3] = true

	ret[POSITION_J1][POSITION_K3] = true
	ret[POSITION_L1][POSITION_K3] = true
	ret[POSITION_J2][POSITION_K3] = true
	ret[POSITION_K2][POSITION_K3] = true
	ret[POSITION_L2][POSITION_K3] = true

	ret[POSITION_N1][POSITION_O3] = true
	ret[POSITION_P1][POSITION_O3] = true
	ret[POSITION_N2][POSITION_O3] = true
	ret[POSITION_O2][POSITION_O3] = true
	ret[POSITION_P2][POSITION_O3] = true

	/* rest of the board from bottom to top */
	ret[POSITION_A3][POSITION_A4] = true
	ret[POSITION_A3][POSITION_B3] = true
	ret[POSITION_B3][POSITION_A3] = true
	ret[POSITION_B3][POSITION_C3] = true
	ret[POSITION_C3][POSITION_B3] = true
	ret[POSITION_C3][POSITION_D3] = true
	ret[POSITION_D3][POSITION_C3] = true
	ret[POSITION_D3][POSITION_E3] = true
	ret[POSITION_E3][POSITION_D3] = true
	ret[POSITION_E3][POSITION_F3] = true
	ret[POSITION_E3][POSITION_E4] = true
	ret[POSITION_F3][POSITION_E3] = true
	ret[POSITION_F3][POSITION_G3] = true
	ret[POSITION_G3][POSITION_F3] = true
	ret[POSITION_G3][POSITION_H3] = true
	ret[POSITION_H3][POSITION_G3] = true
	ret[POSITION_H3][POSITION_I3] = true
	ret[POSITION_I3][POSITION_H3] = true
	ret[POSITION_I3][POSITION_J3] = true
	ret[POSITION_I3][POSITION_I4] = true
	ret[POSITION_J3][POSITION_I3] = true
	ret[POSITION_J3][POSITION_K3] = true
	ret[POSITION_K3][POSITION_J3] = true
	ret[POSITION_K3][POSITION_L3] = true
	ret[POSITION_L3][POSITION_K3] = true
	ret[POSITION_L3][POSITION_M3] = true
	ret[POSITION_M3][POSITION_L3] = true
	ret[POSITION_M3][POSITION_N3] = true
	ret[POSITION_M3][POSITION_M4] = true
	ret[POSITION_N3][POSITION_M3] = true
	ret[POSITION_N3][POSITION_O3] = true
	ret[POSITION_O3][POSITION_N3] = true
	ret[POSITION_O3][POSITION_P3] = true
	ret[POSITION_P3][POSITION_O3] = true
	ret[POSITION_P3][POSITION_Q3] = true
	ret[POSITION_Q3][POSITION_P3] = true
	ret[POSITION_Q3][POSITION_Q4] = true
	ret[POSITION_A4][POSITION_A3] = true
	ret[POSITION_A4][POSITION_A5] = true
	ret[POSITION_E4][POSITION_E3] = true
	ret[POSITION_E4][POSITION_E5] = true
	ret[POSITION_I4][POSITION_I3] = true
	ret[POSITION_I4][POSITION_I5] = true
	ret[POSITION_M4][POSITION_M3] = true
	ret[POSITION_M4][POSITION_M5] = true
	ret[POSITION_Q4][POSITION_Q3] = true
	ret[POSITION_Q4][POSITION_Q5] = true
	ret[POSITION_A5][POSITION_A4] = true
	ret[POSITION_A5][POSITION_B5] = true
	ret[POSITION_B5][POSITION_A5] = true
	ret[POSITION_B5][POSITION_C5] = true
	ret[POSITION_C5][POSITION_B5] = true
	ret[POSITION_C5][POSITION_D5] = true
	ret[POSITION_C5][POSITION_C6] = true
	ret[POSITION_D5][POSITION_C5] = true
	ret[POSITION_D5][POSITION_E5] = true
	ret[POSITION_E5][POSITION_D5] = true
	ret[POSITION_E5][POSITION_F5] = true
	ret[POSITION_E5][POSITION_E4] = true
	ret[POSITION_F5][POSITION_E5] = true
	ret[POSITION_F5][POSITION_G5] = true
	ret[POSITION_G5][POSITION_F5] = true
	ret[POSITION_G5][POSITION_H5] = true
	ret[POSITION_H5][POSITION_G5] = true
	ret[POSITION_H5][POSITION_I5] = true
	ret[POSITION_I5][POSITION_H5] = true
	ret[POSITION_I5][POSITION_J5] = true
	ret[POSITION_I5][POSITION_I4] = true
	ret[POSITION_J5][POSITION_I5] = true
	ret[POSITION_J5][POSITION_K5] = true
	ret[POSITION_K5][POSITION_J5] = true
	ret[POSITION_K5][POSITION_L5] = true
	ret[POSITION_L5][POSITION_K5] = true
	ret[POSITION_L5][POSITION_M5] = true
	ret[POSITION_M5][POSITION_L5] = true
	ret[POSITION_M5][POSITION_N5] = true
	ret[POSITION_M5][POSITION_M4] = true
	ret[POSITION_N5][POSITION_M5] = true
	ret[POSITION_N5][POSITION_O5] = true
	ret[POSITION_O5][POSITION_N5] = true
	ret[POSITION_O5][POSITION_P5] = true
	ret[POSITION_O5][POSITION_O6] = true
	ret[POSITION_P5][POSITION_O5] = true
	ret[POSITION_P5][POSITION_Q5] = true
	ret[POSITION_Q5][POSITION_P5] = true
	ret[POSITION_Q5][POSITION_Q4] = true
	ret[POSITION_C6][POSITION_C5] = true
	ret[POSITION_C6][POSITION_C7] = true
	ret[POSITION_O6][POSITION_O5] = true
	ret[POSITION_O6][POSITION_O7] = true
	ret[POSITION_C7][POSITION_C6] = true
	ret[POSITION_C7][POSITION_D7] = true
	ret[POSITION_D7][POSITION_C7] = true
	ret[POSITION_D7][POSITION_E7] = true
	ret[POSITION_E7][POSITION_D7] = true
	ret[POSITION_E7][POSITION_F7] = true
	ret[POSITION_E7][POSITION_E8] = true
	ret[POSITION_F7][POSITION_E7] = true
	ret[POSITION_F7][POSITION_G7] = true
	ret[POSITION_G7][POSITION_F7] = true
	ret[POSITION_G7][POSITION_H7] = true
	ret[POSITION_H7][POSITION_G7] = true
	ret[POSITION_H7][POSITION_I7] = true
	ret[POSITION_I7][POSITION_H7] = true
	ret[POSITION_I7][POSITION_J7] = true
	ret[POSITION_J7][POSITION_I7] = true
	ret[POSITION_J7][POSITION_K7] = true
	ret[POSITION_K7][POSITION_J7] = true
	ret[POSITION_K7][POSITION_L7] = true
	ret[POSITION_L7][POSITION_K7] = true
	ret[POSITION_L7][POSITION_M7] = true
	ret[POSITION_M7][POSITION_L7] = true
	ret[POSITION_M7][POSITION_N7] = true
	ret[POSITION_M7][POSITION_M8] = true
	ret[POSITION_N7][POSITION_M7] = true
	ret[POSITION_N7][POSITION_O7] = true
	ret[POSITION_O7][POSITION_N7] = true
	ret[POSITION_O7][POSITION_O6] = true
	ret[POSITION_E8][POSITION_E7] = true
	ret[POSITION_E8][POSITION_E9] = true
	ret[POSITION_M8][POSITION_M7] = true
	ret[POSITION_M8][POSITION_M9] = true
	ret[POSITION_E9][POSITION_E8] = true
	ret[POSITION_E9][POSITION_F9] = true
	ret[POSITION_F9][POSITION_E9] = true
	ret[POSITION_F9][POSITION_G9] = true
	ret[POSITION_G9][POSITION_F9] = true
	ret[POSITION_G9][POSITION_H9] = true
	ret[POSITION_G9][POSITION_G10] = true
	ret[POSITION_H9][POSITION_G9] = true
	ret[POSITION_H9][POSITION_I9] = true
	ret[POSITION_I9][POSITION_H9] = true
	ret[POSITION_I9][POSITION_J9] = true
	ret[POSITION_J9][POSITION_I9] = true
	ret[POSITION_J9][POSITION_K9] = true
	ret[POSITION_K9][POSITION_J9] = true
	ret[POSITION_K9][POSITION_L9] = true
	ret[POSITION_K9][POSITION_K10] = true
	ret[POSITION_L9][POSITION_K9] = true
	ret[POSITION_L9][POSITION_M9] = true
	ret[POSITION_M9][POSITION_L9] = true
	ret[POSITION_M9][POSITION_M8] = true
	ret[POSITION_G10][POSITION_G9] = true
	ret[POSITION_G10][POSITION_G11] = true
	ret[POSITION_K10][POSITION_K9] = true
	ret[POSITION_K10][POSITION_K11] = true
	ret[POSITION_G11][POSITION_G10] = true
	ret[POSITION_G11][POSITION_H11] = true
	ret[POSITION_H11][POSITION_G11] = true
	ret[POSITION_H11][POSITION_I11] = true
	ret[POSITION_I11][POSITION_H11] = true
	ret[POSITION_I11][POSITION_J11] = true
	ret[POSITION_I11][POSITION_I12] = true
	ret[POSITION_J11][POSITION_I11] = true
	ret[POSITION_J11][POSITION_K11] = true
	ret[POSITION_K11][POSITION_J11] = true
	ret[POSITION_K11][POSITION_K10] = true
	ret[POSITION_I12][POSITION_I11] = true
	ret[POSITION_I12][POSITION_I13] = true
	ret[POSITION_A13][POSITION_A14] = true
	ret[POSITION_A13][POSITION_B13] = true
	ret[POSITION_B13][POSITION_A13] = true
	ret[POSITION_B13][POSITION_C13] = true
	ret[POSITION_C13][POSITION_B13] = true
	ret[POSITION_C13][POSITION_D13] = true
	ret[POSITION_D13][POSITION_C13] = true
	ret[POSITION_D13][POSITION_E13] = true
	ret[POSITION_E13][POSITION_D13] = true
	ret[POSITION_E13][POSITION_F13] = true
	ret[POSITION_F13][POSITION_E13] = true
	ret[POSITION_F13][POSITION_G13] = true
	ret[POSITION_G13][POSITION_F13] = true
	ret[POSITION_G13][POSITION_H13] = true
	ret[POSITION_H13][POSITION_G13] = true
	ret[POSITION_H13][POSITION_I13] = true
	ret[POSITION_I13][POSITION_H13] = true
	ret[POSITION_I13][POSITION_J13] = true
	ret[POSITION_I13][POSITION_I12] = true
	ret[POSITION_J13][POSITION_I13] = true
	ret[POSITION_J13][POSITION_K13] = true
	ret[POSITION_K13][POSITION_J13] = true
	ret[POSITION_K13][POSITION_L13] = true
	ret[POSITION_L13][POSITION_K13] = true
	ret[POSITION_L13][POSITION_M13] = true
	ret[POSITION_M13][POSITION_L13] = true
	ret[POSITION_M13][POSITION_N13] = true
	ret[POSITION_N13][POSITION_M13] = true
	ret[POSITION_N13][POSITION_O13] = true
	ret[POSITION_O13][POSITION_N13] = true
	ret[POSITION_O13][POSITION_P13] = true
	ret[POSITION_P13][POSITION_O13] = true
	ret[POSITION_P13][POSITION_Q13] = true
	ret[POSITION_Q13][POSITION_P13] = true
	ret[POSITION_Q13][POSITION_Q14] = true
	ret[POSITION_A14][POSITION_A13] = true
	ret[POSITION_A14][POSITION_A15] = true
	ret[POSITION_Q14][POSITION_Q13] = true
	ret[POSITION_Q14][POSITION_Q15] = true
	ret[POSITION_A15][POSITION_A14] = true
	ret[POSITION_A15][POSITION_B15] = true
	ret[POSITION_B15][POSITION_A15] = true
	ret[POSITION_B15][POSITION_C15] = true
	ret[POSITION_C15][POSITION_B15] = true
	ret[POSITION_C15][POSITION_D15] = true
	ret[POSITION_D15][POSITION_C15] = true
	ret[POSITION_D15][POSITION_E15] = true
	ret[POSITION_E15][POSITION_D15] = true
	ret[POSITION_E15][POSITION_F15] = true
	ret[POSITION_F15][POSITION_E15] = true
	ret[POSITION_F15][POSITION_G15] = true
	ret[POSITION_G15][POSITION_F15] = true
	ret[POSITION_G15][POSITION_H15] = true
	ret[POSITION_H15][POSITION_G15] = true
	ret[POSITION_H15][POSITION_I15] = true
	ret[POSITION_I15][POSITION_H15] = true
	ret[POSITION_I15][POSITION_J15] = true
	ret[POSITION_I15][POSITION_I16] = true
	ret[POSITION_J15][POSITION_I15] = true
	ret[POSITION_J15][POSITION_K15] = true
	ret[POSITION_K15][POSITION_J15] = true
	ret[POSITION_K15][POSITION_L15] = true
	ret[POSITION_L15][POSITION_K15] = true
	ret[POSITION_L15][POSITION_M15] = true
	ret[POSITION_M15][POSITION_L15] = true
	ret[POSITION_M15][POSITION_N15] = true
	ret[POSITION_N15][POSITION_M15] = true
	ret[POSITION_N15][POSITION_O15] = true
	ret[POSITION_O15][POSITION_N15] = true
	ret[POSITION_O15][POSITION_P15] = true
	ret[POSITION_P15][POSITION_O15] = true
	ret[POSITION_P15][POSITION_Q15] = true
	ret[POSITION_Q15][POSITION_P15] = true
	ret[POSITION_Q15][POSITION_Q14] = true

	return ret
}()

const (
	STATE_WAITING_FOR_PLAYERS = uint8(0)
	STATE_RED_PLAYING         = uint8(1)
	STATE_GREEN_PLAYING       = uint8(2)
	STATE_YELLOW_PLAYING      = uint8(3)
	STATE_BLUE_PLAYING        = uint8(4)
	STATE_RED_WON             = uint8(5)
	STATE_GREEN_WON           = uint8(6)
	STATE_YELLOW_WON          = uint8(7)
	STATE_BLUE_WON            = uint8(8)
)

var indexStates = []string{
	"WAITING_FOR_PLAYERS",
	"RED_PLAYING",
	"GREEN_PLAYING",
	"YELLOW_PLAYING",
	"BLUE_PLAYING",
	"RED_WON",
	"GREEN_WON",
	"YELLOW_WON",
	"BLUE_WON",
}

var initialPositions = [31]uint8{
	/* RED pawns */
	POSITION_B1, // idx: 0
	POSITION_D1, // idx: 1
	POSITION_B2, // idx: 2
	POSITION_C2, // idx: 3
	POSITION_D2, // idx: 4
	/* GREEN pawns */
	POSITION_F1, // idx: 5
	POSITION_H1, // idx: 6
	POSITION_F2, // idx: 7
	POSITION_G2, // idx: 8
	POSITION_H2, // idx: 9
	/* YELLOW pawns */
	POSITION_J1, // idx: 10
	POSITION_L1, // idx: 11
	POSITION_J2, // idx: 12
	POSITION_K2, // idx: 13
	POSITION_L2, // idx: 14
	/* BLUE pawns */
	POSITION_N1, // idx: 15
	POSITION_P1, // idx: 16
	POSITION_N2, // idx: 17
	POSITION_O2, // idx: 18
	POSITION_P2, // idx: 19
	/* BARICADE pawns */
	POSITION_A5,  // idx: 20
	POSITION_E5,  // idx: 21
	POSITION_I5,  // idx: 22
	POSITION_M5,  // idx: 23
	POSITION_Q5,  // idx: 24
	POSITION_G9,  // idx: 25
	POSITION_K9,  // idx: 26
	POSITION_I11, // idx: 27
	POSITION_I12, // idx: 28
	POSITION_I13, // idx: 29
	POSITION_I15, // idx: 30
}

func InitGame(creator string) *Game {
	return &Game{
		players:   []string{creator},
		state:     STATE_WAITING_FOR_PLAYERS,
		dice:      0,
		positions: initialPositions,
	}
}

func (g *Game) Join(name string) error {
	/* TODO check player not alreayd joined */

	if len(g.players) == 4 {
		return errors.New("game is full")
	}

	g.players = append(g.players, name)

	g.Notify()

	return nil
}

func (g *Game) Start(user string) error {
	if len(g.players) < 2 {
		return errors.New("not enough players to start")
	}

	if user != g.players[0] {
		return errors.New("only game creator can start the game")
	}

	if len(g.players) == 2 {
		/* we double the player's list */
		g.players = append(g.players, g.players...)
	}

	g.state = STATE_RED_PLAYING

	g.Notify()

	return nil
}

func (g *Game) RollDice(user string) error {
	if g.state < STATE_RED_PLAYING {
		return errors.New("game has not started yet")
	}
	if g.state >= STATE_RED_WON {
		return errors.New("game is already complete")
	}

	/* check player is allowed to roll */
	if (g.state == STATE_RED_PLAYING && g.players[0] == user) ||
		(g.state == STATE_GREEN_PLAYING && g.players[1] == user) ||
		(g.state == STATE_YELLOW_PLAYING && g.players[2] == user) ||
		(g.state == STATE_BLUE_PLAYING && g.players[3] == user) {
		rand.Seed(time.Now().UnixNano())
		g.dice = uint8(rand.Intn(6) + 1)

		g.Notify()

		return nil
	}

	return errors.New("not your turn to roll")
}

type path struct {
	visited map[uint8]bool
	path    []uint8
}

func (p *path) append(position uint8) path {
	visited := make(map[uint8]bool)

	for k, v := range p.visited {
		visited[k] = v
	}
	visited[p.path[len(p.path)-1]] = true

	newPath := make([]uint8, len(p.path)+1)
	for i := range p.path {
		newPath[i] = p.path[i]
	}
	newPath[len(newPath)-1] = position

	return path{
		visited: visited,
		path:    newPath,
	}
}

func (p *path) nextPaths() []path {
	/* pre allocate to two because it is the maximum out transitions and not a
	* big space anyway */
	next := make([]path, 0, 2)

	/* get all adjacents positions */
	currentPosition := p.path[len(p.path)-1]
	for n, t := range transitions[currentPosition] {
		if t {
			_, found := p.visited[uint8(n)]
			if !found {
				next = append(next, p.append(uint8(n)))
			}
		}
	}

	return next
}

func findAllPaths(start uint8, dice uint8) [][]uint8 {
	/* initial state */
	paths := []path{path{
		visited: make(map[uint8]bool),
		path:    []uint8{start},
	}}

	for i := uint8(0); len(paths) > 0 && i < dice; i++ {
		newPaths := make([]path, 0, len(paths)*2)

		for i := range paths {
			newPaths = append(newPaths, paths[i].nextPaths()...)
		}

		paths = newPaths
	}

	resultPaths := make([][]uint8, 0, len(paths))
	for i := range paths {
		resultPaths = append(resultPaths, paths[i].path)
	}

	return resultPaths
}

func (g *Game) Move(player, from, to, baricade string) error {
	if g.state < STATE_RED_PLAYING {
		return errors.New("game has not started yet")
	}

	/* check it is player's turn */
	if (g.state == STATE_RED_PLAYING && g.players[0] != player) ||
		(g.state == STATE_GREEN_PLAYING && g.players[1] != player) ||
		(g.state == STATE_YELLOW_PLAYING && g.players[2] != player) ||
		(g.state == STATE_BLUE_PLAYING && g.players[3] != player) {
		return errors.New("not your turn to move")
	}

	if g.dice == 0 {
		return errors.New("roll the dice first")
	}

	playerPawn := ^uint8(0)
	switch g.state {
	case STATE_RED_PLAYING:
		playerPawn = PAWN_RED
		if len(g.players) == 4 && g.players[0] == g.players[2] {
			playerPawn |= PAWN_YELLOW
		}
	case STATE_GREEN_PLAYING:
		playerPawn = PAWN_GREEN
		if len(g.players) == 4 && g.players[1] == g.players[3] {
			playerPawn |= PAWN_BLUE
		}
	case STATE_YELLOW_PLAYING:
		playerPawn = PAWN_YELLOW
		if len(g.players) == 4 && g.players[0] == g.players[2] {
			playerPawn |= PAWN_RED
		}
	case STATE_BLUE_PLAYING:
		playerPawn = PAWN_BLUE
		if len(g.players) == 4 && g.players[1] == g.players[3] {
			playerPawn |= PAWN_GREEN
		}
	}
	if playerPawn == ^uint8(0) {
		return errors.New("invalid action in current game state")
	}

	if to == baricade {
		return errors.New("not allowed baricade destination")
	}

	fromPosition, found := positionToIndex[from]
	if !found {
		return errors.New("invalid from position")
	}
	toPosition, found := positionToIndex[to]
	if !found {
		return errors.New("invalid to position")
	}
	/* baricade is not always required, only when to contains a baricade */
	baricadePosition, _ := positionToIndex[baricade]

	fmt.Println("from, to, baricade:", from, to, baricade)
	fmt.Println("pop, pop, pop:", fromPosition, toPosition, baricadePosition)

	fromPawn, toPawn, baricadePawn := -1, -1, -1
	for i := range g.positions {
		switch g.positions[i] {
		case fromPosition:
			fromPawn = i
		case toPosition:
			toPawn = i
		case baricadePosition:
			baricadePawn = i
		}
	}
	/* check from is player's pawn */
	if fromPawn == -1 {
		return errors.New("start position is empty")
	}
	if (fromPawn < 5 && (playerPawn&PAWN_RED) == 0) ||
		(fromPawn > 4 && fromPawn < 10 && (playerPawn&PAWN_GREEN) == 0) ||
		(fromPawn > 9 && fromPawn < 15 && (playerPawn&PAWN_YELLOW) == 0) ||
		(fromPawn > 14 && fromPawn < 20 && (playerPawn&PAWN_BLUE) == 0) {
		return errors.New("start position is not your own")
	}
	/* check to is allowed destination */
	if (allowedPawns[toPosition] & playerPawn) == 0 {
		return errors.New("not allowed pawn destination")
	}
	isDestinationPlayer := false
	isDestinationBaricade := false
	if toPawn != -1 {
		isDestinationPlayer = toPawn < 20
		isDestinationBaricade = toPawn >= 20
		/* if to is baricade, check baricade is allowed destination and is empty */
		if isDestinationBaricade && baricadePawn != -1 {
			return errors.New("baricade destination is not empty")
		}
		if isDestinationBaricade && (allowedPawns[baricadePosition]&PAWN_BARICADE) == 0 {
			return errors.New("not allowed baricade destination")
		}
	}
	/* generate all possible paths */
	allPaths := findAllPaths(fromPosition, g.dice)
	fmt.Println("allPaths:", allPaths)
	/* keep only paths with to as final element */
	baricadePositions := make(map[uint8]bool)
	for i := 20; i < 31; i++ {
		baricadePositions[g.positions[i]] = true
	}
	i := 0
	for _, p := range allPaths {
		if len(p) != int(g.dice)+1 {
			panic("invalid path length")
		}
		if p[len(p)-1] == toPosition {
			/* check there is no baricade on the path (except for last position) */
			hasBaricade := false
			for j := 0; j < len(p)-1; j++ {
				_, hasBaricade = baricadePositions[p[j]]
				if hasBaricade {
					break
				}
			}

			if !hasBaricade {
				allPaths[i] = p
				i++
			}
		}
	}
	allPaths = allPaths[:i]
	/* check we still have paths in the result set */
	if len(allPaths) == 0 {
		return errors.New("no available path between source and destination (check baricades)")
	}
	/* do the move */
	g.positions[fromPawn] = toPosition
	if isDestinationPlayer {
		g.positions[toPawn] = initialPositions[toPawn]
	} else if isDestinationBaricade {
		g.positions[toPawn] = baricadePosition
	}

	/* check if game is finished */
	if toPosition == POSITION_I16 {
		g.state += uint8(4)
	} else {
		g.state -= STATE_RED_PLAYING
		g.state += uint8(1)
		g.state %= uint8(len(g.players))
		g.state += STATE_RED_PLAYING
	}

	/* reset dice */
	g.dice = 0

	/* notify update to connected players */
	g.Notify()

	/* TODO return full path */
	return nil
}

type GameJson struct {
	Players   []string          `json:"players"`
	State     string            `json:"state"`
	Dice      int               `json:"dice"`
	Positions map[string]string `json:"positions"`
}

func (g *Game) JSON() *GameJson {
	/* generate positions in JSON */
	positionsToPawn := make(map[uint8]string)
	for i, p := range g.positions {
		pawn := "baricade"
		if i < 5 {
			pawn = "red"
		} else if i > 4 && i < 10 {
			pawn = "green"
		} else if i > 9 && i < 15 {
			pawn = "yellow"
		} else if i > 14 && i < 20 {
			pawn = "blue"
		}

		positionsToPawn[uint8(p)] = pawn
	}
	positions := make(map[string]string)
	for k, p := range positionToIndex {
		pawn, found := positionsToPawn[p]
		if !found {
			pawn = ""
		}
		positions[strings.ToLower(k)] = pawn
	}

	return &GameJson{
		Players:   g.players,
		State:     indexStates[g.state],
		Dice:      int(g.dice),
		Positions: positions,
	}
}

func (g *Game) Subscribe() chan *GameJson {
	newChan := make(chan *GameJson, 4)

	g.subscriptions = append(g.subscriptions, newChan)

	return newChan
}

func (g *Game) Notify() {
	update := g.JSON()
	for _, c := range g.subscriptions {
		c <- update
	}
}

func (g *Game) Unsubscribe(toDelete chan *GameJson) {
	for i, c := range g.subscriptions {
		if c == toDelete {
			g.subscriptions = append(g.subscriptions[:i], g.subscriptions[i+1:]...)
			return
		}
	}
}
