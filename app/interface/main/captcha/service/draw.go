package service

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"math/rand"
	"time"

	"go-common/app/interface/main/captcha/conf"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// CONST VALUE.
const (
	NORMAL      = int(4)
	MEDIUM      = int(8)
	HIGH        = int(16)
	MinLenStart = int(4)
	MinWidth    = int(48)
	MinLength   = int(20)
	Length48    = int(48)

	TypeNone  = int(0)
	TypeLOWER = int(1)
	TypeUPPER = int(2)
	TypeALL   = int(3)
)

var fontKinds = [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}

func sign(x int) int {
	if x > 0 {
		return 1
	}
	return -1
}

// NewCaptcha new a captcha.
func newCaptcha(c *conf.Captcha) *Captcha {
	captcha := &Captcha{
		disturbLevel: NORMAL,
	}
	captcha.frontColors = []color.Color{color.Black}
	captcha.bkgColors = []color.Color{color.White}
	captcha.setFont(c.Fonts...)
	colors := []color.Color{}
	for _, v := range c.BkgColors {
		colors = append(colors, v)
	}
	captcha.setBkgColor(colors...)
	colors = []color.Color{}
	for _, v := range c.FrontColors {
		colors = append(colors, v)
	}
	captcha.setFontColor(colors...)
	captcha.setDisturbance(c.DisturbLevel)
	return captcha
}

// addFont add font.
func (c *Captcha) addFont(path string) error {
	fontdata, erro := ioutil.ReadFile(path)
	if erro != nil {
		return erro
	}
	font, erro := freetype.ParseFont(fontdata)
	if erro != nil {
		return erro
	}
	if c.fonts == nil {
		c.fonts = []*truetype.Font{}
	}
	c.fonts = append(c.fonts, font)
	return nil
}

// setFont set font.
func (c *Captcha) setFont(paths ...string) (err error) {
	for _, v := range paths {
		if err = c.addFont(v); err != nil {
			return err
		}
	}
	return nil
}

// setBkgColor set backgroud color.
func (c *Captcha) setBkgColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.bkgColors = c.bkgColors[:0]
		c.bkgColors = append(c.bkgColors, colors...)
	}
}

func (c *Captcha) randFont() *truetype.Font {
	return c.fonts[rand.Intn(len(c.fonts))]
}

// setBkgsetFontColorColor set font color.
func (c *Captcha) setFontColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.frontColors = c.frontColors[:0]
		c.frontColors = append(c.frontColors, colors...)
	}
}

// setDisturbance set disturbance.
func (c *Captcha) setDisturbance(d int) {
	if d > 0 {
		c.disturbLevel = d
	}
}

func (c *Captcha) createImage(lenStart, lenEnd, width, length, t int) (image *Image, str string) {
	num := MinLenStart
	if lenStart < MinLenStart {
		lenStart = MinLenStart
	}
	if lenEnd > lenStart {
		// rand.Seed(time.Now().UnixNano())
		num = rand.Intn(lenEnd-lenStart+1) + lenStart
	}
	str = c.randStr(num, t)
	return c.createCustom(str, width, length), str
}

func (c *Captcha) createCustom(str string, width, length int) *Image {
	// boundary check
	if len(str) == 0 {
		str = "bilibili"
	}
	if width < MinWidth {
		width = MinWidth
	}
	if length < MinLength {
		length = MinLength
	}
	dst := newImage(width, length)
	c.drawBkg(dst)
	c.drawNoises(dst)
	c.drawString(dst, str, width, length)
	return dst
}

// randStr ascII random
// 48~57 -> 0~9 number
// 65~90 -> A~Z uppercase
// 98~122 -> a~z lowcase
func (c *Captcha) randStr(size, kind int) string {
	ikind, result := kind, make([]byte, size)
	isAll := kind > TypeUPPER || kind < TypeNone
	// rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll {
			ikind = rand.Intn(TypeALL)
		}
		scope, base := fontKinds[ikind][0], fontKinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

func (c *Captcha) drawBkg(img *Image) {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	//??????????????????
	bgcolorindex := ra.Intn(len(c.bkgColors))
	bkg := image.NewUniform(c.bkgColors[bgcolorindex])
	img.fillBkg(bkg)
}

func (c *Captcha) drawNoises(img *Image) {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	//// ????????????????????????
	point := img.Bounds().Size()
	disturbLevel := c.disturbLevel
	// ??????????????????
	for i := 0; i < disturbLevel; i++ {
		x := ra.Intn(point.X)
		y := ra.Intn(point.Y)
		radius := ra.Intn(point.Y/20) + 1
		colorindex := ra.Intn(len(c.frontColors))
		img.drawCircle(x, y, radius, i%4 != 0, c.frontColors[colorindex])
	}
	// ???????????????
	for i := 0; i < disturbLevel; i++ {
		x := ra.Intn(point.X)
		y := ra.Intn(point.Y)
		o := int(math.Pow(-1, float64(i)))
		w := ra.Intn(point.Y) * o
		h := ra.Intn(point.Y/10) * o
		colorindex := ra.Intn(len(c.frontColors))
		img.drawLine(x, y, x+w, y+h, c.frontColors[colorindex])
		colorindex++
	}
}

// ????????????
func (c *Captcha) drawString(img *Image, str string, width, length int) {

	if c.fonts == nil {
		panic("????????????????????????")
	}
	tmp := newImage(width, length)

	// ?????????????????????????????? 0.6
	fsize := int(float64(length) * 0.6)
	// ????????????????????????
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// ?????????????????????
	// ?????????????????????1/4?????????????????????
	padding := fsize / 4
	gap := (width - padding*2) / (len(str))

	// ??????????????????????????????
	for i, char := range str {
		// ????????????????????????
		// ??????????????????????????????????????????
		str := newImage(fsize, fsize)
		// str.FillBkg(image.NewUniform(color.Black))
		// ????????????????????????
		colorindex := r.Intn(len(c.frontColors))

		//?????????????????????
		font := c.randFont()
		str.drawString(font, c.frontColors[colorindex], string(char), float64(fsize))

		// ??????????????????????????????
		rs := str.rotate(float64(r.Intn(40) - 20))
		// ??????????????????
		s := rs.Bounds().Size()
		left := i*gap + padding
		top := (length - s.Y) / 2
		// ??????????????????
		draw.Draw(tmp, image.Rect(left, top, left+s.X, top+s.Y), rs, image.ZP, draw.Over)
	}
	if length >= Length48 {
		// ????????????48???????????? ??????48????????????????????????
		tmp.distortTo(float64(fsize)/10, 200.0)
	}
	draw.Draw(img, tmp.Bounds(), tmp, image.ZP, draw.Over)
}
