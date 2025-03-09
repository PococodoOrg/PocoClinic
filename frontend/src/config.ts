interface Config {
  apiUrl: string
  environment: string
  version: string
  auth: {
    tokenKey: string
    refreshTokenKey: string
  }
  pagination: {
    defaultPageSize: number
    maxPageSize: number
  }
}

const defaultConfig: Config = {
  apiUrl: 'http://localhost:8080',
  environment: 'development',
  version: '1.0.0',
  auth: {
    tokenKey: 'poco_token',
    refreshTokenKey: 'poco_refresh_token'
  },
  pagination: {
    defaultPageSize: 10,
    maxPageSize: 100
  }
}

const envConfigs: Record<string, Partial<Config>> = {
  development: {
    apiUrl: 'http://localhost:8080'
  },
  test: {
    apiUrl: 'http://localhost:8080'
  },
  production: {
    apiUrl: 'https://api.pococlinic.com'
  }
}

const env = import.meta.env.MODE || 'development'
const envConfig = envConfigs[env] || {}

export const config: Config = {
  ...defaultConfig,
  ...envConfig,
  environment: env
}

// Type-safe config getter
export function getConfig<K extends keyof Config>(key: K): Config[K] {
  return config[key]
} 