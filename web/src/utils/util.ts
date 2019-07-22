/**
 * hex to rgba
 * @param {string} hex Color hex value
 * @param {number} alpha
 * @return {string} rgba
 */
export function toRGBA(hex: string, alpha: number) {
  const hexReg = /^#?([0-9a-fA-f]{3}|[0-9a-fA-f]{6})$/

  if (hexReg.test(hex)) {
    const validAlpha = Math.max(0, Math.min(1, alpha))
    const validHex = hex.startsWith('#') ? hex.slice(1) : hex

    const length = validHex.length === 3 ? 1 : 2

    const color = []
    for (let i = 0; i < validHex.length; i += length) {
      color.push(parseInt(`0x${validHex.slice(i, i + length)}`))
    }

    return `rgba(${color.join(', ')}, ${validAlpha})`
  } else {
    return hex
  }
}

/**
 * rgb to hex
 * @param red
 * @param green
 * @param blue
 * @return {string} Color
 */
export function toHex(red: string | number, green: string | number, blue: string | number) {
  return `#${red.toString(16).padStart(0)}${green.toString(16).padStart(0)}${blue.toString(16).padStart(0)}`
}

/**
 * Get chart series color
 * @param {number} idx index
 * @return {string} Color hex value
 */
export function getChartColor(idx: number): string {
  const colors = [ '#7EB26D', '#EAB839', '#6ED0E0', '#EF843C', '#E24D42', '#1F78C1', '#BA43A9', '#705DA0', '#508642', '#CCA300', '#447EBC', '#C15C17', '#890F02', '#0A437C', '#6D1F62', '#584477', '#70DBED', '#F9BA8F', '#F29191', '#82B5D8', '#E5A8E2', '#AEA2E0', '#629E51', '#E5AC0E', '#64B0C8', '#E0752D', '#BF1B00', '#0A50A1', '#962D82', '#614D93', '#9AC48A', '#F2C96D', '#65C5DB', '#F9934E', '#EA6460', '#5195CE', '#D683CE', '#806EB7', '#3F6833', '#967302', '#2F575E', '#99440A', '#58140C', '#052B51', '#511749', '#3F2B5B' ] // tslint:disable-line
  return colors[ idx % colors.length ]
}