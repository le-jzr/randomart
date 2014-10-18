package randomart

type AugmentFunc func(x, y int)
type StepFunc func(x, y, width, heigth, inst int) (nextx, nexty int)

const (
	SSH_FLDSIZE_X = 17
	SSH_FLDSIZE_Y = 9
)

func DiagonalStep(x, y, maxx, maxy, inst int) (nextx, nexty int) {
	if (inst & 0x1) != 0 {
		if x+1 < maxx {
			x++
		}
	} else {
		if x > 0 {
			x--
		}
	}

	if (inst & 0x2) != 0 {
		if y+1 < maxy {
			y++
		}
	} else {
		if y > 0 {
			y--
		}
	}

	return x, y
}

func GridWrapStep(x, y, maxx, maxy, inst int) (nextx, nexty int) {
	switch inst {
	case 0:
		x++
	case 1:
		x--
	case 2:
		y++
	case 3:
		y--
	}

	if x < 0 {
		x = maxx - 1
	}
	if y < 0 {
		y = maxy - 1
	}

	return x % maxx, y % maxy
}

func OctogonalStep(x, y, maxx, maxy, inst int) (nextx, nexty int) {
	switch inst {
	case 0:
		x++
	case 1:
		x++
		y++
	case 2:
		y++
	case 3:
		x--
		y++
	case 4:
		x--
	case 5:
		x--
		y--
	case 6:
		y--
	case 7:
		x++
		y--
	}

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= maxx {
		x = maxx - 1
	}
	if y >= maxy {
		y = maxy - 1
	}

	return x, y
}

func OpenSSH(instructions []byte) (ret [SSH_FLDSIZE_Y][SSH_FLDSIZE_X]byte) {
	const augmentation_string = " .o+=*BOX@%&#/^SE"

	var field [SSH_FLDSIZE_X][SSH_FLDSIZE_Y]int

	augment := func(x, y int) {
		field[x][y]++
	}

	Generic(instructions, 2, SSH_FLDSIZE_X/2, SSH_FLDSIZE_Y/2, SSH_FLDSIZE_X, SSH_FLDSIZE_Y, DiagonalStep, augment)

	for x := 0; x < SSH_FLDSIZE_X; x++ {
		for y := 0; y < SSH_FLDSIZE_Y; y++ {
			val := field[x][y]
			if val > len(augmentation_string)-1 {
				val = len(augmentation_string) - 1
			}
			ret[y][x] = augmentation_string[val]
		}
	}

	return
}

func Generic(instructions []byte, isize uint, startx, starty, width, height int, step StepFunc, augment AugmentFunc) {
	xpos := startx
	ypos := starty

	var register uint64
	var registerBits uint

	for _, b := range instructions {

		register |= (uint64(b) << registerBits)
		registerBits += 8

		for registerBits >= isize {

			inst := int(register & ((1 << isize) - 1))
			register >>= isize
			registerBits -= isize

			xpos, ypos = step(xpos, ypos, width, height, inst)
			augment(xpos, ypos)
		}
	}
}
