import { http, HttpResponse } from 'msw'
import type {
  LoginRequest,
  LoginResponse,
  User,
  Patient,
  PaginatedResponse,
  ApiError
} from '../types/api'

const BASE_URL = 'http://localhost:8080'

// Mock data
const mockUser: User = {
  id: '1',
  name: 'Test User',
  email: 'test@example.com',
  role: 'doctor',
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString()
}

const mockPatients: Patient[] = [
  {
    id: '1',
    firstName: 'John',
    lastName: 'Doe',
    dateOfBirth: '1990-01-01',
    gender: 'male',
    email: 'john.doe@example.com',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }
]

// Error responses
const createErrorResponse = (error: ApiError, status = 400) => {
  return HttpResponse.json(error, { status })
}

export const handlers = [
  // Mock the login endpoint
  http.post<never, LoginRequest>(`${BASE_URL}/api/auth/login`, async ({ request }) => {
    try {
      const body = await request.json() as LoginRequest
      if (!body?.email || !body?.password) {
        return createErrorResponse({
          message: 'Email and password are required',
          details: {
            email: !body?.email ? ['Email is required'] : [],
            password: !body?.password ? ['Password is required'] : []
          }
        })
      }

      const response: LoginResponse = {
        token: 'mock-token',
        user: mockUser
      }

      return HttpResponse.json(response)
    } catch (error) {
      return createErrorResponse({
        message: 'Internal server error',
        code: 'INTERNAL_ERROR'
      }, 500)
    }
  }),

  // Mock the current user endpoint
  http.get(`${BASE_URL}/api/auth/me`, ({ request }) => {
    const authHeader = request.headers.get('Authorization')
    if (!authHeader?.startsWith('Bearer ')) {
      return createErrorResponse({
        message: 'Unauthorized',
        code: 'UNAUTHORIZED'
      }, 401)
    }
    return HttpResponse.json(mockUser)
  }),

  // Mock the patients list endpoint
  http.get(`${BASE_URL}/api/v1/patients`, ({ request }) => {
    try {
      const url = new URL(request.url)
      const page = parseInt(url.searchParams.get('page') || '1')
      const pageSize = parseInt(url.searchParams.get('pageSize') || '10')
      const search = url.searchParams.get('search') || ''

      // Validate query parameters
      if (isNaN(page) || page < 1 || isNaN(pageSize) || pageSize < 1) {
        return createErrorResponse({
          message: 'Invalid pagination parameters',
          details: {
            page: isNaN(page) || page < 1 ? ['Page must be a positive number'] : [],
            pageSize: isNaN(pageSize) || pageSize < 1 ? ['Page size must be a positive number'] : []
          }
        })
      }

      // Filter patients by search term
      const filteredPatients = search
        ? mockPatients.filter(p => 
            `${p.firstName} ${p.lastName}`.toLowerCase().includes(search.toLowerCase()) ||
            p.email?.toLowerCase().includes(search.toLowerCase())
          )
        : mockPatients

      const response: PaginatedResponse<Patient> = {
        items: filteredPatients,
        total: filteredPatients.length,
        page,
        pageSize
      }

      return HttpResponse.json(response)
    } catch (error) {
      return createErrorResponse({
        message: 'Internal server error',
        code: 'INTERNAL_ERROR'
      }, 500)
    }
  }),

  // Add more handlers as needed
] 