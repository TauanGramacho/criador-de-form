export type ClientOptions = {
    baseUrl: `${string}://${string}` | (string & {});
};

export type AuthConfigResponse = {

    readonly $schema?: string;
    googleEnabled: boolean;
};

export type AuthResponse = {

    readonly $schema?: string;
    user: User;
};

export type CreateFormRequest = {

    readonly $schema?: string;
    description?: string;
    fields: Array<FormField> | null;
    title: string;
};

export type ErrorDetail = {

    location?: string;

    message?: string;

    value?: unknown;
};

export type ErrorModel = {

    readonly $schema?: string;

    detail?: string;

    errors?: Array<ErrorDetail> | null;

    instance?: string;

    status?: number;

    title?: string;

    type?: string;
};

export type FormDetail = {

    readonly $schema?: string;
    createdAt: string;
    description: string;
    fields: Array<FormField> | null;
    id: string;
    publicUrl?: string;
    publishedAt?: string;
    responseCount: number;
    slug: string;
    status: 'draft' | 'published';
    title: string;
    updatedAt: string;
};

export type FormField = {
    helpText?: string;
    id: string;
    label: string;
    max?: number;
    min?: number;
    options?: Array<string> | null;
    placeholder?: string;
    required: boolean;
    type: 'text' | 'textarea' | 'email' | 'number' | 'select' | 'checkbox';
};

export type FormResponse = {
    answers: {
        [key: string]: unknown;
    };
    formId: string;
    id: string;
    submittedAt: string;
};

export type FormSummary = {
    createdAt: string;
    description: string;
    fieldCount: number;
    id: string;
    publicUrl?: string;
    publishedAt?: string;
    responseCount: number;
    slug: string;
    status: 'draft' | 'published';
    title: string;
    updatedAt: string;
};

export type LoginRequest = {

    readonly $schema?: string;
    email: string;
    password: string;
};

export type PublicForm = {

    readonly $schema?: string;
    description: string;
    fields: Array<FormField> | null;
    id: string;
    slug: string;
    title: string;
};

export type RegisterRequest = {

    readonly $schema?: string;
    confirmPassword: string;
    email: string;
    name: string;
    password: string;
};

export type SubmitResponseRequest = {

    readonly $schema?: string;
    answers: {
        [key: string]: unknown;
    };
};

export type SubmitResponseResult = {

    readonly $schema?: string;
    id: string;
    submittedAt: string;
};

export type User = {
    createdAt: string;
    email: string;
    id: string;
    name: string;
};

export type AuthConfigResponseWritable = {
    googleEnabled: boolean;
};

export type AuthResponseWritable = {
    user: User;
};

export type CreateFormRequestWritable = {
    description?: string;
    fields: Array<FormField> | null;
    title: string;
};

export type ErrorModelWritable = {

    detail?: string;

    errors?: Array<ErrorDetail> | null;

    instance?: string;

    status?: number;

    title?: string;

    type?: string;
};

export type FormDetailWritable = {
    createdAt: string;
    description: string;
    fields: Array<FormField> | null;
    id: string;
    publicUrl?: string;
    publishedAt?: string;
    responseCount: number;
    slug: string;
    status: 'draft' | 'published';
    title: string;
    updatedAt: string;
};

export type LoginRequestWritable = {
    email: string;
    password: string;
};

export type PublicFormWritable = {
    description: string;
    fields: Array<FormField> | null;
    id: string;
    slug: string;
    title: string;
};

export type RegisterRequestWritable = {
    confirmPassword: string;
    email: string;
    name: string;
    password: string;
};

export type SubmitResponseRequestWritable = {
    answers: {
        [key: string]: unknown;
    };
};

export type SubmitResponseResultWritable = {
    id: string;
    submittedAt: string;
};

export type GetAuthConfigData = {
    body?: never;
    path?: never;
    query?: never;
    url: '/api/auth/config';
};

export type GetAuthConfigErrors = {

    default: ErrorModel;
};

export type GetAuthConfigError = GetAuthConfigErrors[keyof GetAuthConfigErrors];

export type GetAuthConfigResponses = {

    200: AuthConfigResponse;
};

export type GetAuthConfigResponse = GetAuthConfigResponses[keyof GetAuthConfigResponses];

export type CompleteGoogleLoginData = {
    body?: never;
    path?: never;
    query?: {
        code?: string;
        state?: string;
    };
    url: '/api/auth/google/callback';
};

export type CompleteGoogleLoginErrors = {

    default: ErrorModel;
};

export type CompleteGoogleLoginError = CompleteGoogleLoginErrors[keyof CompleteGoogleLoginErrors];

export type CompleteGoogleLoginResponses = {

    default: ErrorModel;
};

export type CompleteGoogleLoginResponse = CompleteGoogleLoginResponses[keyof CompleteGoogleLoginResponses];

export type StartGoogleLoginData = {
    body?: never;
    path?: never;
    query?: {
        redirect?: string;
    };
    url: '/api/auth/google/start';
};

export type StartGoogleLoginErrors = {

    default: ErrorModel;
};

export type StartGoogleLoginError = StartGoogleLoginErrors[keyof StartGoogleLoginErrors];

export type StartGoogleLoginResponses = {

    default: ErrorModel;
};

export type StartGoogleLoginResponse = StartGoogleLoginResponses[keyof StartGoogleLoginResponses];

export type LoginWithEmailData = {
    body: LoginRequestWritable;
    path?: never;
    query?: never;
    url: '/api/auth/login';
};

export type LoginWithEmailErrors = {

    default: ErrorModel;
};

export type LoginWithEmailError = LoginWithEmailErrors[keyof LoginWithEmailErrors];

export type LoginWithEmailResponses = {

    200: AuthResponse;
};

export type LoginWithEmailResponse = LoginWithEmailResponses[keyof LoginWithEmailResponses];

export type LogoutData = {
    body?: never;
    path?: never;
    query?: never;
    url: '/api/auth/logout';
};

export type LogoutErrors = {

    default: ErrorModel;
};

export type LogoutError = LogoutErrors[keyof LogoutErrors];

export type LogoutResponses = {

    204: void;
};

export type LogoutResponse = LogoutResponses[keyof LogoutResponses];

export type GetCurrentUserData = {
    body?: never;
    path?: never;
    query?: never;
    url: '/api/auth/me';
};

export type GetCurrentUserErrors = {

    default: ErrorModel;
};

export type GetCurrentUserError = GetCurrentUserErrors[keyof GetCurrentUserErrors];

export type GetCurrentUserResponses = {

    200: AuthResponse;
};

export type GetCurrentUserResponse = GetCurrentUserResponses[keyof GetCurrentUserResponses];

export type RegisterWithEmailData = {
    body: RegisterRequestWritable;
    path?: never;
    query?: never;
    url: '/api/auth/register';
};

export type RegisterWithEmailErrors = {

    default: ErrorModel;
};

export type RegisterWithEmailError = RegisterWithEmailErrors[keyof RegisterWithEmailErrors];

export type RegisterWithEmailResponses = {

    201: AuthResponse;
};

export type RegisterWithEmailResponse = RegisterWithEmailResponses[keyof RegisterWithEmailResponses];

export type ListFormsData = {
    body?: never;
    path?: never;
    query?: never;
    url: '/api/forms';
};

export type ListFormsErrors = {

    default: ErrorModel;
};

export type ListFormsError = ListFormsErrors[keyof ListFormsErrors];

export type ListFormsResponses = {

    200: Array<FormSummary> | null;
};

export type ListFormsResponse = ListFormsResponses[keyof ListFormsResponses];

export type CreateFormData = {
    body: CreateFormRequestWritable;
    path?: never;
    query?: never;
    url: '/api/forms';
};

export type CreateFormErrors = {

    default: ErrorModel;
};

export type CreateFormError = CreateFormErrors[keyof CreateFormErrors];

export type CreateFormResponses = {

    201: FormDetail;
};

export type CreateFormResponse = CreateFormResponses[keyof CreateFormResponses];

export type DeleteFormData = {
    body?: never;
    path: {
        id: string;
    };
    query?: never;
    url: '/api/forms/{id}';
};

export type DeleteFormErrors = {

    default: ErrorModel;
};

export type DeleteFormError = DeleteFormErrors[keyof DeleteFormErrors];

export type DeleteFormResponses = {

    204: void;
};

export type DeleteFormResponse = DeleteFormResponses[keyof DeleteFormResponses];

export type GetFormData = {
    body?: never;
    path: {
        id: string;
    };
    query?: never;
    url: '/api/forms/{id}';
};

export type GetFormErrors = {

    default: ErrorModel;
};

export type GetFormError = GetFormErrors[keyof GetFormErrors];

export type GetFormResponses = {

    200: FormDetail;
};

export type GetFormResponse = GetFormResponses[keyof GetFormResponses];

export type UpdateFormData = {
    body: CreateFormRequestWritable;
    path: {
        id: string;
    };
    query?: never;
    url: '/api/forms/{id}';
};

export type UpdateFormErrors = {

    default: ErrorModel;
};

export type UpdateFormError = UpdateFormErrors[keyof UpdateFormErrors];

export type UpdateFormResponses = {

    200: FormDetail;
};

export type UpdateFormResponse = UpdateFormResponses[keyof UpdateFormResponses];

export type PublishFormData = {
    body?: never;
    path: {
        id: string;
    };
    query?: never;
    url: '/api/forms/{id}/publish';
};

export type PublishFormErrors = {

    default: ErrorModel;
};

export type PublishFormError = PublishFormErrors[keyof PublishFormErrors];

export type PublishFormResponses = {

    200: FormDetail;
};

export type PublishFormResponse = PublishFormResponses[keyof PublishFormResponses];

export type ListFormResponsesData = {
    body?: never;
    path: {
        id: string;
    };
    query?: never;
    url: '/api/forms/{id}/responses';
};

export type ListFormResponsesErrors = {

    default: ErrorModel;
};

export type ListFormResponsesError = ListFormResponsesErrors[keyof ListFormResponsesErrors];

export type ListFormResponsesResponses = {

    200: Array<FormResponse> | null;
};

export type ListFormResponsesResponse = ListFormResponsesResponses[keyof ListFormResponsesResponses];

export type UnpublishFormData = {
    body?: never;
    path: {
        id: string;
    };
    query?: never;
    url: '/api/forms/{id}/unpublish';
};

export type UnpublishFormErrors = {

    default: ErrorModel;
};

export type UnpublishFormError = UnpublishFormErrors[keyof UnpublishFormErrors];

export type UnpublishFormResponses = {

    200: FormDetail;
};

export type UnpublishFormResponse = UnpublishFormResponses[keyof UnpublishFormResponses];

export type GetPublicFormData = {
    body?: never;
    path: {
        slug: string;
    };
    query?: never;
    url: '/api/public/forms/{slug}';
};

export type GetPublicFormErrors = {

    default: ErrorModel;
};

export type GetPublicFormError = GetPublicFormErrors[keyof GetPublicFormErrors];

export type GetPublicFormResponses = {

    200: PublicForm;
};

export type GetPublicFormResponse = GetPublicFormResponses[keyof GetPublicFormResponses];

export type SubmitPublicFormResponseData = {
    body: SubmitResponseRequestWritable;
    headers?: {
        'X-Forwarded-For'?: string;
    };
    path: {
        slug: string;
    };
    query?: never;
    url: '/api/public/forms/{slug}/responses';
};

export type SubmitPublicFormResponseErrors = {

    default: ErrorModel;
};

export type SubmitPublicFormResponseError = SubmitPublicFormResponseErrors[keyof SubmitPublicFormResponseErrors];

export type SubmitPublicFormResponseResponses = {

    201: SubmitResponseResult;
};

export type SubmitPublicFormResponseResponse = SubmitPublicFormResponseResponses[keyof SubmitPublicFormResponseResponses];
