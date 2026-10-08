package internal

const MAX_CART_SIZE = 1024 * 1024 // 1MB
var cartridge_data [MAX_CART_SIZE]uint8
var cartridge_loaded bool = false
