export interface User {
  id: number;
  email: string;
}

export interface AuthResponse {
  message: string;
  token?: string;
  user: User;
}

export interface LoginCredentials {
  email: string;
  password?: string;
}
