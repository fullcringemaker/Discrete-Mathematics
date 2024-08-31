package main

func encode(utf32 []rune) []byte {
    var utf8 []byte
    for _, r := range utf32 {
        switch {
        case r <= 0x7F:
            utf8 = append(utf8, byte(r))
        case r <= 0x7FF:
            utf8 = append(utf8, byte(0xC0|(r>>6)))
            utf8 = append(utf8, byte(0x80|(r&0x3F)))
        case r <= 0xFFFF:
            utf8 = append(utf8, byte(0xE0|(r>>12)))
            utf8 = append(utf8, byte(0x80|((r>>6)&0x3F)))
            utf8 = append(utf8, byte(0x80|(r&0x3F)))
        default:
            utf8 = append(utf8, byte(0xF0|(r>>18)))
            utf8 = append(utf8, byte(0x80|((r>>12)&0x3F)))
            utf8 = append(utf8, byte(0x80|((r>>6)&0x3F)))
            utf8 = append(utf8, byte(0x80|(r&0x3F)))
        }
    }
    return utf8
}

func decode(utf8 []byte) []rune {
    var utf32 []rune
    var r rune
    var size byte
    for _, b := range utf8 {
        switch {
        case b&0x80 == 0x00:
            utf32 = append(utf32, rune(b))
        case b&0xE0 == 0xC0:
            r = rune(b & 0x1F)
            size = 1
        case b&0xF0 == 0xE0:
            r = rune(b & 0x0F)
            size = 2
        case b&0xF8 == 0xF0:
            r = rune(b & 0x07)
            size = 3
        case b&0xC0 == 0x80:
            r = (r << 6) | rune(b&0x3F)
            size--
            if size == 0 {
                utf32 = append(utf32, r)
            }
        }
    }
    return utf32
}

func main() {

}
