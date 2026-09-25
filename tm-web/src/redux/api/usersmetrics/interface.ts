export interface WeeklyPercentageChange {
  active_users_percentage_change: number;
  average_time_percentage_change: number;
  new_users_percentage_change: number;
  total_sessions_percentage_change: number;
}

export interface ActiveAndNewUsers {
  active_users: number;
  new_users: number;
}

export interface MetricsData {
  average_time_per_session: number;
  total_active_users: number;
  total_sessions: number;
  total_users: number;
  daily_active_and_new_users: ActiveAndNewUsers;
  weekly_active_and_new_users: ActiveAndNewUsers;
  monthly_active_and_new_users: ActiveAndNewUsers;
  weekly_percentage_change: WeeklyPercentageChange;
}

export interface IUsersMetrics {
  message: string;
  data: MetricsData;
  timestamp: string;
  status: string;
}
export interface IWalletDistributionResponse {
  data: string;
  errors?: string;
  message: string;
  status: string;
  timestamp: string;
}

export interface Referrer {
  referrer: string;
  count: number;
}

export interface IReferrerCountResponse {
  message: string;
  status: string;
  timestamp: string;
  data: Referrer[];
}

export interface RecentRegistration {
  id: string;
  username: string;
  email: string;
  address: string;
  phone_number: string;
  full_name: string;
  date: string;
}

export interface RecentRegistrationsResponse {
  data: {
    data: RecentRegistration[];
    total: number;
    page: number;
    pageSize: number;
  };
  message: string;
  status: string;
  timestamp: string;
}

export interface ICountriesStats {
  data: {
    Country: string;
    Count: number;
  }[];
  message: string;
  status: string;
  timestamp: string;
}
