import type { Client, Options as Options2, TDataShape } from './client';
import { client } from './client.gen';
import type { CompleteGoogleLoginData, CompleteGoogleLoginErrors, CompleteGoogleLoginResponses, CreateFormData, CreateFormErrors, CreateFormResponses, DeleteFormData, DeleteFormErrors, DeleteFormResponses, GetAuthConfigData, GetAuthConfigErrors, GetAuthConfigResponses, GetCurrentUserData, GetCurrentUserErrors, GetCurrentUserResponses, GetFormData, GetFormErrors, GetFormResponses, GetPublicFormData, GetPublicFormErrors, GetPublicFormResponses, ListFormResponsesData, ListFormResponsesErrors, ListFormResponsesResponses, ListFormsData, ListFormsErrors, ListFormsResponses, LoginWithEmailData, LoginWithEmailErrors, LoginWithEmailResponses, LogoutData, LogoutErrors, LogoutResponses, PublishFormData, PublishFormErrors, PublishFormResponses, RegisterWithEmailData, RegisterWithEmailErrors, RegisterWithEmailResponses, StartGoogleLoginData, StartGoogleLoginErrors, StartGoogleLoginResponses, SubmitPublicFormResponseData, SubmitPublicFormResponseErrors, SubmitPublicFormResponseResponses, UnpublishFormData, UnpublishFormErrors, UnpublishFormResponses, UpdateFormData, UpdateFormErrors, UpdateFormResponses } from './types.gen';

export type Options<TData extends TDataShape = TDataShape, ThrowOnError extends boolean = boolean, TResponse = unknown> = Options2<TData, ThrowOnError, TResponse> & {

    client?: Client;

    meta?: Record<string, unknown>;
};


export const getAuthConfig = <ThrowOnError extends boolean = false>(options?: Options<GetAuthConfigData, ThrowOnError>) => (options?.client ?? client).get<GetAuthConfigResponses, GetAuthConfigErrors, ThrowOnError>({ url: '/api/auth/config', ...options });


export const completeGoogleLogin = <ThrowOnError extends boolean = false>(options?: Options<CompleteGoogleLoginData, ThrowOnError>) => (options?.client ?? client).get<CompleteGoogleLoginResponses, CompleteGoogleLoginErrors, ThrowOnError>({ url: '/api/auth/google/callback', ...options });


export const startGoogleLogin = <ThrowOnError extends boolean = false>(options?: Options<StartGoogleLoginData, ThrowOnError>) => (options?.client ?? client).get<StartGoogleLoginResponses, StartGoogleLoginErrors, ThrowOnError>({ url: '/api/auth/google/start', ...options });


export const loginWithEmail = <ThrowOnError extends boolean = false>(options: Options<LoginWithEmailData, ThrowOnError>) => (options.client ?? client).post<LoginWithEmailResponses, LoginWithEmailErrors, ThrowOnError>({
    url: '/api/auth/login',
    ...options,
    headers: {
        'Content-Type': 'application/json',
        ...options.headers
    }
});


export const logout = <ThrowOnError extends boolean = false>(options?: Options<LogoutData, ThrowOnError>) => (options?.client ?? client).post<LogoutResponses, LogoutErrors, ThrowOnError>({ url: '/api/auth/logout', ...options });


export const getCurrentUser = <ThrowOnError extends boolean = false>(options?: Options<GetCurrentUserData, ThrowOnError>) => (options?.client ?? client).get<GetCurrentUserResponses, GetCurrentUserErrors, ThrowOnError>({ url: '/api/auth/me', ...options });


export const registerWithEmail = <ThrowOnError extends boolean = false>(options: Options<RegisterWithEmailData, ThrowOnError>) => (options.client ?? client).post<RegisterWithEmailResponses, RegisterWithEmailErrors, ThrowOnError>({
    url: '/api/auth/register',
    ...options,
    headers: {
        'Content-Type': 'application/json',
        ...options.headers
    }
});


export const listForms = <ThrowOnError extends boolean = false>(options?: Options<ListFormsData, ThrowOnError>) => (options?.client ?? client).get<ListFormsResponses, ListFormsErrors, ThrowOnError>({ url: '/api/forms', ...options });


export const createForm = <ThrowOnError extends boolean = false>(options: Options<CreateFormData, ThrowOnError>) => (options.client ?? client).post<CreateFormResponses, CreateFormErrors, ThrowOnError>({
    url: '/api/forms',
    ...options,
    headers: {
        'Content-Type': 'application/json',
        ...options.headers
    }
});


export const deleteForm = <ThrowOnError extends boolean = false>(options: Options<DeleteFormData, ThrowOnError>) => (options.client ?? client).delete<DeleteFormResponses, DeleteFormErrors, ThrowOnError>({ url: '/api/forms/{id}', ...options });


export const getForm = <ThrowOnError extends boolean = false>(options: Options<GetFormData, ThrowOnError>) => (options.client ?? client).get<GetFormResponses, GetFormErrors, ThrowOnError>({ url: '/api/forms/{id}', ...options });


export const updateForm = <ThrowOnError extends boolean = false>(options: Options<UpdateFormData, ThrowOnError>) => (options.client ?? client).put<UpdateFormResponses, UpdateFormErrors, ThrowOnError>({
    url: '/api/forms/{id}',
    ...options,
    headers: {
        'Content-Type': 'application/json',
        ...options.headers
    }
});


export const publishForm = <ThrowOnError extends boolean = false>(options: Options<PublishFormData, ThrowOnError>) => (options.client ?? client).post<PublishFormResponses, PublishFormErrors, ThrowOnError>({ url: '/api/forms/{id}/publish', ...options });


export const listFormResponses = <ThrowOnError extends boolean = false>(options: Options<ListFormResponsesData, ThrowOnError>) => (options.client ?? client).get<ListFormResponsesResponses, ListFormResponsesErrors, ThrowOnError>({ url: '/api/forms/{id}/responses', ...options });


export const unpublishForm = <ThrowOnError extends boolean = false>(options: Options<UnpublishFormData, ThrowOnError>) => (options.client ?? client).post<UnpublishFormResponses, UnpublishFormErrors, ThrowOnError>({ url: '/api/forms/{id}/unpublish', ...options });


export const getPublicForm = <ThrowOnError extends boolean = false>(options: Options<GetPublicFormData, ThrowOnError>) => (options.client ?? client).get<GetPublicFormResponses, GetPublicFormErrors, ThrowOnError>({ url: '/api/public/forms/{slug}', ...options });


export const submitPublicFormResponse = <ThrowOnError extends boolean = false>(options: Options<SubmitPublicFormResponseData, ThrowOnError>) => (options.client ?? client).post<SubmitPublicFormResponseResponses, SubmitPublicFormResponseErrors, ThrowOnError>({
    url: '/api/public/forms/{slug}/responses',
    ...options,
    headers: {
        'Content-Type': 'application/json',
        ...options.headers
    }
});
