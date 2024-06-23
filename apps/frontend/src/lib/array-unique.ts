export function arrayUnique<T>(array: T[]) {
  return Array.from(new Set(array))
}

export function arrayUniqueBy<T>(
  array: T[],
  fn: (a: T, b: T, array: T[]) => boolean
) {
  return array.reduce<T[]>((acc, v) => {
    if (!acc.some((x) => fn(v, x, array))) acc.push(v)
    return acc
  }, [])
}
