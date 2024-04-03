import type { Config } from 'jest'

const config: Config = {
  verbose: true,
  testPathIgnorePatterns: ['/node_modules/', '/dist/']
}

export default config
