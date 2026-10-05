import { baseApi } from './baseapi';
import { tagTypes } from './baseapi/tagTypes';
import { P2PCreds } from './p2pApis';
import { PMAsset, PMDividend, PMOrder, PMPortfolio, PMPriceRange, PMPrices, PMQuote } from '../../types/publicMarkets';

type WithCreds<T> = T & { creds: P2PCreds };

export type PMTradeBody = {
  walletAddress: string;
  amount?: string; // buy: the funding asset to spend
  quantity?: string; // sell: the tokens to sell
  transaction?: string;
  transactionSignature?: string;
};

// The first call (no signature) returns the operation to sign; the second,
// signed one places the order.
export type PMTradeResponse = { quote: PMQuote; transaction?: string; messages?: string[]; order?: PMOrder };

const qs = (params: Record<string, string | undefined>) => {
  const s = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => v && s.set(k, v));
  const out = s.toString();
  return out ? `?${out}` : '';
};
const assetUrl = (code: string) => `/v1/public-markets/assets/${encodeURIComponent(code)}`;

// publicMarketsApi: tokenized NGX equities and FMDQ bonds. The catalogue,
// asset detail and prices are public (unsigned); quotes, trades, the
// portfolio, orders and dividends are signed with the user's credentials.
// A trade is two calls, as with bank withdrawals: the unsigned build, then
// the same call with the signed transaction.
export const publicMarketsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    listPMAssets: builder.query<{ assets: PMAsset[] }, { market?: string; type?: string; search?: string }>({
      query: (p) => ({ url: `/v1/public-markets${qs(p)}`, method: 'GET' }),
    }),
    getPMAsset: builder.query<PMAsset, string>({
      query: (code) => ({ url: assetUrl(code), method: 'GET' }),
    }),
    getPMPrices: builder.query<PMPrices, { code: string; range: PMPriceRange }>({
      query: ({ code, range }) => ({ url: `${assetUrl(code)}/prices${qs({ range })}`, method: 'GET' }),
    }),
    getPMQuote: builder.query<PMQuote, WithCreds<{ code: string; side: 'buy' | 'sell'; value: string }>>({
      query: ({ creds, code, side, value }) => ({
        url: `${assetUrl(code)}/quote${qs(side === 'buy' ? { side, amount: value } : { side, quantity: value })}`,
        method: 'GET',
        data: { creds },
      }),
    }),
    tradePM: builder.mutation<PMTradeResponse, WithCreds<{ code: string; side: 'buy' | 'sell'; body: PMTradeBody }>>({
      query: ({ creds, code, side, body }) => ({ url: `${assetUrl(code)}/${side}`, method: 'POST', data: { payload: body, creds } }),
      invalidatesTags: (r) => (r?.order ? [tagTypes.publicMarkets] : []),
    }),
    getPMPortfolio: builder.query<PMPortfolio, WithCreds<{}>>({
      query: ({ creds }) => ({ url: '/v1/public-markets/portfolio', method: 'GET', data: { creds } }),
      providesTags: [tagTypes.publicMarkets],
    }),
    getPMOrders: builder.query<{ orders: PMOrder[] }, WithCreds<{}>>({
      query: ({ creds }) => ({ url: '/v1/public-markets/orders', method: 'GET', data: { creds } }),
      providesTags: [tagTypes.publicMarkets],
    }),
    getPMOrder: builder.query<PMOrder, WithCreds<{ id: string }>>({
      query: ({ creds, id }) => ({ url: `/v1/public-markets/orders/${encodeURIComponent(id)}`, method: 'GET', data: { creds } }),
    }),
    getPMDividends: builder.query<{ dividends: PMDividend[] }, WithCreds<{ assetCode?: string }>>({
      query: ({ creds, assetCode }) => ({ url: `/v1/public-markets/dividends${qs({ assetCode })}`, method: 'GET', data: { creds } }),
      providesTags: [tagTypes.publicMarkets],
    }),
  }),
});

export const {
  useListPMAssetsQuery,
  useGetPMAssetQuery,
  useGetPMPricesQuery,
  useGetPMQuoteQuery,
  useTradePMMutation,
  useGetPMPortfolioQuery,
  useGetPMOrdersQuery,
  useGetPMOrderQuery,
  useGetPMDividendsQuery,
} = publicMarketsApi;
