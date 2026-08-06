export const BASE_URL = import.meta.env.VITE_API_URL ?? '/api';

const wsHost = typeof window !== 'undefined' ? window.location.host : '';
export const SOCKET_URL = import.meta.env.VITE_SOCKET_URL ?? `wss://${wsHost}/ws/v1`;
