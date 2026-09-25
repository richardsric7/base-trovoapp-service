interface ISuccessResponse<T = void> {
  success: boolean;
  message: string;
  data: T;
}

interface IPaginatedResponse<T> {
  status: number;
  message: string;
  data: {
    data: T[];
    page: number;
    pageSize: number;
    total: number;
    hasNextPage: boolean;
    hasPreviousPage: boolean;
  };
}
interface IErrorResponse {
  success: boolean;
  status: number;
  message: string;
  error: string;
  errors: string[];
}
