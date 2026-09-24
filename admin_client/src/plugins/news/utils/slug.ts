const randomPart = () => {
  const values = crypto.getRandomValues(new Uint32Array(2))
  return Array.from(values, value => value.toString(36)).join('')
}

// createRandomSlug 生成仅包含小写字母和数字的业务标识。
export const createRandomSlug = (prefix: 'category' | 'article') => `${prefix}${Date.now().toString(36)}${randomPart()}`
