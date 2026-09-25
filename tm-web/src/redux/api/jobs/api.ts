import { baseApi } from "@/redux/baseApi";
import {
  GetJobsParams,
  GetJobsResponse,
  JobPostPayload,
  JobResponse,
} from "./interface";

export const jobApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getJobs: builder.query<GetJobsResponse, GetJobsParams>({
      query: (params) => ({
        url: "/admin/careers/roles",
        method: "GET",
        params,
      }),
      providesTags: ["jobs"],
    }),
    getJobById: builder.query({
      query: (id) => ({
        url: `/admin/careers/roles/${id}`,
        method: "GET",
      }),
      providesTags: ["jobs"],
    }),
    postJob: builder.mutation<JobResponse, JobPostPayload>({
      query: (data) => ({
        url: "/admin/careers/roles",
        method: "POST",
        data: data,
      }),
      invalidatesTags: ["jobs"],
    }),
    publishJob: builder.mutation({
      query: (id) => ({
        url: `/admin/careers/roles/${id}/publish`,
        method: "PUT",
      }),
      invalidatesTags: ["jobs"],
    }),

    jobStatus: builder.mutation({
      query: (params: { id: string; status: "published" | "disabled" }) => ({
        url: `/admin/careers/roles/${params.id}/status`,
        method: "PUT",
        data: { status: params.status },
      }),
      invalidatesTags: ["jobs"],
    }),
    updateJob: builder.mutation({
      query: (data) => ({
        url: `/admin/careers/roles/${data.id}`,
        method: "PUT",
        data: data,
      }),
      invalidatesTags: ["jobs"],
    }),
    deleteJob: builder.mutation({
      query: (id) => ({
        url: `/admin/careers/roles/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["jobs"],
    }),
  }),
});

export const {
  useGetJobsQuery,
  useGetJobByIdQuery,
  usePostJobMutation,
  usePublishJobMutation,
  useUpdateJobMutation,
  useDeleteJobMutation,
  useJobStatusMutation,
} = jobApi;
