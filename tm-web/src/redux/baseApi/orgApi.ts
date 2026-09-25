import { createApi } from "@reduxjs/toolkit/query/react";
import { tagTypes } from "./tagTypes";
import { orgAxiosBaseQuery } from "./orgBasedQuery";

export const orgApi = createApi({
  reducerPath: "orgApi",
  baseQuery: orgAxiosBaseQuery({}),
  tagTypes: Object.values(tagTypes),
  endpoints: () => ({}),
});
