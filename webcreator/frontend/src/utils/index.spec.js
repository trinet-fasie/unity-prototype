import { sortObjectsByKey, convertToISODate } from '.'

describe('utils/index.js', () => {
  describe('convertToISODate', () => {
    it('new Date() should be momentjs object', () => {
      expect(convertToISODate(new Date())._isAMomentObject).toBe(true)
    })
  })

  describe('sortObjectsByKey', () => {
    const a = {text: 'a'}
    const b = {text: 'b'}
    const A = {text: 'A'}
    const numberOne = {number: 1}
    const numberTwo = {number: 2}

    it('sort of {text: \'a\'} {text: \'b\'} should return -1', () => {
      expect(sortObjectsByKey(a, b, 'text')).toBe(-1)
    })

    it('sort of {text: \'b\'} {text: \'A\'} should return 1', () => {
      expect(sortObjectsByKey(b, A, 'text')).toBe(1)
    })

    it('sort of {text: \'A\'} {text: \'A\'} should return 0', () => {
      expect(sortObjectsByKey(A, A, 'text')).toBe(0)
    })

    it('sort of {number: 1} {number: 2} should return -1', () => {
      expect(sortObjectsByKey(numberOne, numberTwo, 'number')).toBe(-1)
    })

    it('sort of {number: 2} {number: 1} should return 1', () => {
      expect(sortObjectsByKey(numberTwo, numberOne, 'number')).toBe(1)
    })

    it('sort of {number: 1} {number: 1} should return 0', () => {
      expect(sortObjectsByKey(numberOne, numberOne, 'number')).toBe(0)
    })
  })
})
