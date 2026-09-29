const randomPart = () => {
  const values = crypto.getRandomValues(new Uint32Array(2))
  return Array.from(values, value => value.toString(36)).join('')
}

const slugPrefixes = {
  category: 'c',
  article: 'n',
} as const

// createRandomSlug 统一生成短业务标识：1位类型前缀 + 6位时间片 + 4位安全随机值。
export const createRandomSlug = (type: keyof typeof slugPrefixes) => {
  const timePart = Date.now().toString(36).slice(-6)
  const secureRandomPart = randomPart().slice(0, 4)
  return `${slugPrefixes[type]}${timePart}${secureRandomPart}`
}
