import AddIcon from '@mui/icons-material/Add';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import DeleteIcon from '@mui/icons-material/Delete';
import GoogleIcon from '@mui/icons-material/Google';
import LogoutIcon from '@mui/icons-material/Logout';
import SaveIcon from '@mui/icons-material/Save';
import SendIcon from '@mui/icons-material/Send';
import UnpublishedIcon from '@mui/icons-material/Unpublished';
import {
  Alert,
  AppBar,
  Box,
  Button,
  Checkbox,
  Chip,
  CircularProgress,
  Container,
  FormControl,
  FormControlLabel,
  IconButton,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Stack,
  Tab,
  Tabs,
  TextField,
  Toolbar,
  Tooltip,
  Typography,
} from '@mui/material';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect, useMemo, useState } from 'react';
import {
  BrowserRouter,
  Link as RouterLink,
  Navigate,
  Outlet,
  Route,
  Routes,
  useLocation,
  useNavigate,
  useParams,
} from 'react-router-dom';
import { apiBaseUrl, unwrap } from './api';
import {
  createForm,
  getCurrentUser,
  getAuthConfig,
  getForm,
  getPublicForm,
  listFormResponses,
  listForms,
  loginWithEmail,
  logout,
  publishForm,
  registerWithEmail,
  submitPublicFormResponse,
  unpublishForm,
  updateForm,
} from './api/generated/sdk.gen';
import type { FormDetail, FormField, FormResponse, FormSummary, PublicForm, User } from './api/generated/types.gen';

type AuthMode = 'login' | 'register';
type FieldType = FormField['type'];
type Answers = Record<string, unknown>;

const fieldTypes: Array<{ value: FieldType; label: string }> = [
  { value: 'text', label: 'Texto curto' },
  { value: 'textarea', label: 'Texto longo' },
  { value: 'email', label: 'E-mail' },
  { value: 'number', label: 'Número' },
  { value: 'select', label: 'Seleção' },
  { value: 'checkbox', label: 'Checkbox' },
];

const defaultFields: FormField[] = [
  {
    id: 'name',
    label: 'Nome',
    placeholder: 'Seu nome',
    required: true,
    type: 'text',
  },
  {
    id: 'email',
    label: 'E-mail',
    placeholder: 'voce@empresa.com',
    required: true,
    type: 'email',
  },
];

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/privacy" element={<LegalPage type="privacy" />} />
        <Route path="/terms" element={<LegalPage type="terms" />} />
        <Route path="/f/:slug" element={<PublicFormPage />} />
        <Route element={<RequireAuth />}>
          <Route path="/" element={<Navigate to="/forms" replace />} />
          <Route path="/forms" element={<DashboardPage />} />
          <Route path="/forms/:id" element={<FormEditorPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

function RequireAuth() {
  const query = useCurrentUser();
  const location = useLocation();

  if (query.isLoading) {
    return <CenteredLoader />;
  }
  if (query.isError || !query.data) {
    return <Navigate to="/login" state={{ from: location.pathname }} replace />;
  }
  return <Shell user={query.data.user} />;
}

function Shell({ user }: { user: User }) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const logoutMutation = useMutation({
    mutationFn: async () => {
      await logout();
    },
    onSuccess: async () => {
      queryClient.clear();
      navigate('/login', { replace: true });
    },
  });

  return (
    <>
      <AppBar position="sticky" color="inherit" elevation={0} sx={{ borderBottom: '1px solid', borderColor: 'divider' }}>
        <Toolbar sx={{ gap: 2 }}>
          <Typography component={RouterLink} to="/forms" variant="h6" color="text.primary" sx={{ textDecoration: 'none', flex: 1 }}>
            Form Builder
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {user.name}
          </Typography>
          <Tooltip title="Sair">
            <IconButton aria-label="Sair" onClick={() => logoutMutation.mutate()} disabled={logoutMutation.isPending}>
              <LogoutIcon />
            </IconButton>
          </Tooltip>
        </Toolbar>
      </AppBar>
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Outlet />
      </Container>
    </>
  );
}

function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const queryClient = useQueryClient();
  const [mode, setMode] = useState<AuthMode>('login');
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  const authConfigQuery = useQuery({
    queryKey: ['auth-config'],
    queryFn: async () => unwrap(await getAuthConfig()),
    retry: false,
  });

  const from = typeof location.state === 'object' && location.state && 'from' in location.state ? String(location.state.from) : '/forms';

  const mutation = useMutation({
    mutationFn: async () => {
      if (mode === 'register') {
        if (password !== confirmPassword) {
          throw new Error('As senhas não conferem.');
        }
        return unwrap(await registerWithEmail({ body: { name, email, password, confirmPassword } }));
      }
      return unwrap(await loginWithEmail({ body: { email, password } }));
    },
    onSuccess: async () => {
      if (mode === 'register') {
        setMode('login');
        setSuccessMessage('Conta criada. Entre com seu e-mail e senha para acessar.');
        setPassword('');
        setConfirmPassword('');
        return;
      }
      await queryClient.invalidateQueries({ queryKey: ['me'] });
      navigate(from, { replace: true });
    },
  });

  return (
    <Box className="auth-page">
      <Paper elevation={0} className="auth-panel">
        <Stack spacing={3}>
          <Box>
            <Typography variant="h4" fontWeight={700}>
              Form Builder
            </Typography>
            <Typography color="text.secondary" sx={{ mt: 1 }}>
              Área administrativa
            </Typography>
          </Box>
          <Tabs
            value={mode}
            onChange={(_, value: AuthMode) => {
              setMode(value);
              setSuccessMessage('');
            }}
          >
            <Tab label="Entrar" value="login" />
            <Tab label="Criar conta" value="register" />
          </Tabs>
          {successMessage ? <Alert severity="success">{successMessage}</Alert> : null}
          {authConfigQuery.data && !authConfigQuery.data.googleEnabled ? (
            <Alert severity="info">Login com Google ainda não está configurado no backend.</Alert>
          ) : null}
          {mutation.error ? <Alert severity="error">{mutation.error.message}</Alert> : null}
          <Box
            component="form"
            onSubmit={(event) => {
              event.preventDefault();
              mutation.mutate();
            }}
          >
            <Stack spacing={2}>
              {mode === 'register' ? (
                <TextField label="Nome" value={name} onChange={(event) => setName(event.target.value)} required fullWidth />
              ) : null}
              <TextField label="E-mail" value={email} onChange={(event) => setEmail(event.target.value)} type="email" required fullWidth />
              <TextField
                label="Senha"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                type="password"
                required
                fullWidth
              />
              {mode === 'register' ? (
                <TextField
                  label="Confirmar senha"
                  value={confirmPassword}
                  onChange={(event) => setConfirmPassword(event.target.value)}
                  type="password"
                  required
                  error={confirmPassword.length > 0 && password !== confirmPassword}
                  helperText={confirmPassword.length > 0 && password !== confirmPassword ? 'As senhas precisam ser iguais.' : ' '}
                  fullWidth
                />
              ) : null}
              <Button type="submit" variant="contained" size="large" disabled={mutation.isPending}>
                {mode === 'register' ? 'Criar conta' : 'Entrar'}
              </Button>
              <Button
                variant="outlined"
                size="large"
                startIcon={<GoogleIcon />}
                disabled={authConfigQuery.isLoading || !authConfigQuery.data?.googleEnabled}
                onClick={() => {
                  window.location.href = `${apiBaseUrl}/api/auth/google/start?redirect=${encodeURIComponent(from)}`;
                }}
              >
                Entrar com Google
              </Button>
            </Stack>
          </Box>
        </Stack>
      </Paper>
    </Box>
  );
}

function DashboardPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const formsQuery = useQuery({
    queryKey: ['forms'],
    queryFn: async () => unwrap(await listForms()) ?? [],
  });

  const createMutation = useMutation({
    mutationFn: async () =>
      unwrap(
        await createForm({
          body: {
            title: 'Novo formulário',
            description: '',
            fields: defaultFields,
          },
        }),
      ),
    onSuccess: async (form) => {
      await queryClient.invalidateQueries({ queryKey: ['forms'] });
      navigate(`/forms/${form.id}`);
    },
  });

  if (formsQuery.isLoading) {
    return <CenteredLoader />;
  }

  const forms = formsQuery.data ?? [];

  return (
    <Stack spacing={3}>
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'center' }}>
        <Box sx={{ flex: 1 }}>
          <Typography variant="h4" fontWeight={700}>
            Formulários
          </Typography>
          <Typography color="text.secondary">Crie, publique e acompanhe respostas.</Typography>
        </Box>
        <Button variant="contained" startIcon={<AddIcon />} onClick={() => createMutation.mutate()} disabled={createMutation.isPending}>
          Novo formulário
        </Button>
      </Stack>
      {formsQuery.error ? <Alert severity="error">{formsQuery.error.message}</Alert> : null}
      {forms.length === 0 ? (
        <Paper elevation={0} className="empty-state">
          <Typography variant="h6">Nenhum formulário ainda</Typography>
          <Button sx={{ mt: 2 }} variant="contained" startIcon={<AddIcon />} onClick={() => createMutation.mutate()}>
            Criar primeiro formulário
          </Button>
        </Paper>
      ) : (
        <Box className="form-grid">
          {forms.map((form) => (
            <FormSummaryCard key={form.id} form={form} />
          ))}
        </Box>
      )}
    </Stack>
  );
}

function FormSummaryCard({ form }: { form: FormSummary }) {
  return (
    <Paper elevation={0} className="summary-card" component={RouterLink} to={`/forms/${form.id}`}>
      <Stack spacing={2}>
        <Stack direction="row" spacing={1} alignItems="center">
          <Chip size="small" color={form.status === 'published' ? 'success' : 'default'} label={form.status === 'published' ? 'Publicado' : 'Rascunho'} />
          <Typography variant="body2" color="text.secondary">
            {form.responseCount} respostas
          </Typography>
        </Stack>
        <Box>
          <Typography variant="h6">{form.title}</Typography>
          <Typography color="text.secondary" className="line-clamp">
            {form.description || 'Sem descrição'}
          </Typography>
        </Box>
        <Typography variant="body2" color="text.secondary">
          {form.fieldCount} campos
        </Typography>
      </Stack>
    </Paper>
  );
}

function FormEditorPage() {
  const { id } = useParams();
  const queryClient = useQueryClient();
  const [tab, setTab] = useState(0);
  const formQuery = useQuery({
    queryKey: ['forms', id],
    queryFn: async () => unwrap(await getForm({ path: { id: id! } })),
    enabled: Boolean(id),
  });
  const responsesQuery = useQuery({
    queryKey: ['forms', id, 'responses'],
    queryFn: async () => unwrap(await listFormResponses({ path: { id: id! } })) ?? [],
    enabled: Boolean(id) && tab === 1,
  });

  if (formQuery.isLoading) {
    return <CenteredLoader />;
  }
  if (formQuery.error || !formQuery.data) {
    return <Alert severity="error">{formQuery.error?.message || 'Formulário não encontrado'}</Alert>;
  }

  return (
    <Stack spacing={3}>
      <FormEditorHeader form={formQuery.data} />
      <Tabs value={tab} onChange={(_, value: number) => setTab(value)}>
        <Tab label="Builder" />
        <Tab label={`Respostas (${formQuery.data.responseCount})`} />
      </Tabs>
      {tab === 0 ? (
        <FormBuilder
          form={formQuery.data}
          onSaved={async () => {
            await queryClient.invalidateQueries({ queryKey: ['forms'] });
            await queryClient.invalidateQueries({ queryKey: ['forms', id] });
          }}
        />
      ) : (
        <ResponsesTable responses={responsesQuery.data ?? []} loading={responsesQuery.isLoading} />
      )}
    </Stack>
  );
}

function FormEditorHeader({ form }: { form: FormDetail }) {
  const queryClient = useQueryClient();
  const publishMutation = useMutation({
    mutationFn: async () =>
      unwrap(
        form.status === 'published'
          ? await unpublishForm({ path: { id: form.id } })
          : await publishForm({ path: { id: form.id } }),
      ),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['forms'] });
      await queryClient.invalidateQueries({ queryKey: ['forms', form.id] });
    },
  });

  return (
    <Paper elevation={0} className="toolbar-panel">
      <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} alignItems={{ md: 'center' }}>
        <Box sx={{ flex: 1 }}>
          <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1 }}>
            <Chip size="small" color={form.status === 'published' ? 'success' : 'default'} label={form.status === 'published' ? 'Publicado' : 'Rascunho'} />
            <Typography variant="body2" color="text.secondary">
              /f/{form.slug}
            </Typography>
          </Stack>
          <Typography variant="h5" fontWeight={700}>
            {form.title}
          </Typography>
        </Box>
        {form.publicUrl ? (
          <Button startIcon={<ContentCopyIcon />} variant="outlined" onClick={() => navigator.clipboard.writeText(form.publicUrl!)}>
            Copiar link
          </Button>
        ) : null}
        <Button
          startIcon={form.status === 'published' ? <UnpublishedIcon /> : <CheckCircleIcon />}
          variant="contained"
          color={form.status === 'published' ? 'inherit' : 'success'}
          onClick={() => publishMutation.mutate()}
          disabled={publishMutation.isPending}
        >
          {form.status === 'published' ? 'Despublicar' : 'Publicar'}
        </Button>
      </Stack>
    </Paper>
  );
}

function FormBuilder({ form, onSaved }: { form: FormDetail; onSaved: () => Promise<void> }) {
  const [title, setTitle] = useState(form.title);
  const [description, setDescription] = useState(form.description);
  const [fields, setFields] = useState<FormField[]>(form.fields ?? []);

  useEffect(() => {
    setTitle(form.title);
    setDescription(form.description);
    setFields(form.fields ?? []);
  }, [form]);

  const saveMutation = useMutation({
    mutationFn: async () =>
      unwrap(
        await updateForm({
          path: { id: form.id },
          body: {
            title,
            description,
            fields,
          },
        }),
      ),
    onSuccess: onSaved,
  });

  const addField = (type: FieldType = 'text') => {
    setFields((current) => [
      ...current,
      {
        id: `field_${Math.random().toString(36).slice(2, 8)}`,
        label: 'Novo campo',
        required: false,
        type,
        options: type === 'select' ? ['Opção 1', 'Opção 2'] : undefined,
      },
    ]);
  };

  const updateField = (fieldID: string, next: Partial<FormField>) => {
    setFields((current) => current.map((field) => (field.id === fieldID ? { ...field, ...next } : field)));
  };

  const removeField = (fieldID: string) => {
    setFields((current) => current.filter((field) => field.id !== fieldID));
  };

  return (
    <Stack spacing={3}>
      {saveMutation.error ? <Alert severity="error">{saveMutation.error.message}</Alert> : null}
      {saveMutation.isSuccess ? <Alert severity="success">Formulário salvo.</Alert> : null}
      <Paper elevation={0} className="editor-panel">
        <Stack spacing={2}>
          <TextField label="Título" value={title} onChange={(event) => setTitle(event.target.value)} required fullWidth />
          <TextField
            label="Descrição"
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            multiline
            minRows={2}
            fullWidth
          />
        </Stack>
      </Paper>
      <Stack spacing={2}>
        {fields.map((field, index) => (
          <Paper key={field.id} elevation={0} className="field-panel">
            <Stack spacing={2}>
              <Stack direction="row" spacing={1} alignItems="center">
                <Typography variant="subtitle2" color="text.secondary" sx={{ flex: 1 }}>
                  Campo {index + 1}
                </Typography>
                <Tooltip title="Remover campo">
                  <IconButton aria-label="Remover campo" onClick={() => removeField(field.id)} disabled={fields.length <= 1}>
                    <DeleteIcon />
                  </IconButton>
                </Tooltip>
              </Stack>
              <Box className="field-grid">
                <TextField label="ID" value={field.id} onChange={(event) => updateField(field.id, { id: normalizeFieldID(event.target.value) })} />
                <TextField label="Rótulo" value={field.label} onChange={(event) => updateField(field.id, { label: event.target.value })} />
                <FormControl>
                  <InputLabel>Tipo</InputLabel>
                  <Select
                    label="Tipo"
                    value={field.type}
                    onChange={(event) => {
                      const type = event.target.value as FieldType;
                      updateField(field.id, {
                        type,
                        options: type === 'select' ? field.options?.length ? field.options : ['Opção 1'] : undefined,
                      });
                    }}
                  >
                    {fieldTypes.map((type) => (
                      <MenuItem key={type.value} value={type.value}>
                        {type.label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
                <FormControlLabel
                  control={<Checkbox checked={field.required} onChange={(event) => updateField(field.id, { required: event.target.checked })} />}
                  label="Obrigatório"
                />
              </Box>
              <Box className="field-grid">
                <TextField
                  label="Placeholder"
                  value={field.placeholder ?? ''}
                  onChange={(event) => updateField(field.id, { placeholder: event.target.value || undefined })}
                />
                <TextField
                  label="Ajuda"
                  value={field.helpText ?? ''}
                  onChange={(event) => updateField(field.id, { helpText: event.target.value || undefined })}
                />
                {field.type === 'number' ? (
                  <>
                    <TextField
                      label="Mínimo"
                      type="number"
                      value={field.min ?? ''}
                      onChange={(event) => updateField(field.id, { min: event.target.value === '' ? undefined : Number(event.target.value) })}
                    />
                    <TextField
                      label="Máximo"
                      type="number"
                      value={field.max ?? ''}
                      onChange={(event) => updateField(field.id, { max: event.target.value === '' ? undefined : Number(event.target.value) })}
                    />
                  </>
                ) : null}
              </Box>
              {field.type === 'select' ? (
                <TextField
                  label="Opções"
                  value={(field.options ?? []).join(', ')}
                  onChange={(event) =>
                    updateField(field.id, {
                      options: event.target.value
                        .split(',')
                        .map((option) => option.trim())
                        .filter(Boolean),
                    })
                  }
                  fullWidth
                />
              ) : null}
            </Stack>
          </Paper>
        ))}
      </Stack>
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
        <Button variant="outlined" startIcon={<AddIcon />} onClick={() => addField()}>
          Adicionar campo
        </Button>
        <Button variant="contained" startIcon={<SaveIcon />} onClick={() => saveMutation.mutate()} disabled={saveMutation.isPending}>
          Salvar alterações
        </Button>
      </Stack>
    </Stack>
  );
}

function ResponsesTable({ responses, loading }: { responses: FormResponse[]; loading: boolean }) {
  if (loading) {
    return <CenteredLoader />;
  }
  if (responses.length === 0) {
    return (
      <Paper elevation={0} className="empty-state">
        <Typography variant="h6">Nenhuma resposta recebida</Typography>
      </Paper>
    );
  }
  return (
    <Stack spacing={2}>
      {responses.map((response) => (
        <Paper key={response.id} elevation={0} className="response-panel">
          <Stack spacing={1}>
            <Typography variant="body2" color="text.secondary">
              {new Date(response.submittedAt).toLocaleString()}
            </Typography>
            {Object.entries(response.answers).map(([key, value]) => (
              <Stack key={key} direction={{ xs: 'column', sm: 'row' }} spacing={1}>
                <Typography variant="body2" fontWeight={700} sx={{ minWidth: 160 }}>
                  {key}
                </Typography>
                <Typography variant="body2">{formatAnswer(value)}</Typography>
              </Stack>
            ))}
          </Stack>
        </Paper>
      ))}
    </Stack>
  );
}

function PublicFormPage() {
  const { slug } = useParams();
  const formQuery = useQuery({
    queryKey: ['public-form', slug],
    queryFn: async () => unwrap(await getPublicForm({ path: { slug: slug! } })),
    enabled: Boolean(slug),
    retry: false,
  });

  if (formQuery.isLoading) {
    return <CenteredLoader />;
  }
  if (formQuery.error || !formQuery.data) {
    return (
      <Container maxWidth="sm" sx={{ py: 8 }}>
        <Alert severity="error">{formQuery.error?.message || 'Formulário indisponível'}</Alert>
      </Container>
    );
  }
  return <PublicForm form={formQuery.data} />;
}

function LegalPage({ type }: { type: 'privacy' | 'terms' }) {
  const isPrivacy = type === 'privacy';

  return (
    <Container maxWidth="md" sx={{ py: 8 }}>
      <Paper elevation={0} className="public-panel">
        <Stack spacing={3}>
          <Box>
            <Typography variant="h4" fontWeight={700}>
              {isPrivacy ? 'Política de Privacidade' : 'Termos de Serviço'}
            </Typography>
            <Typography color="text.secondary" sx={{ mt: 1 }}>
              Form Builder
            </Typography>
          </Box>
          {isPrivacy ? (
            <Stack spacing={2}>
              <Typography>
                O Form Builder usa autenticação do Google apenas para identificar o usuário administrador e permitir acesso à área de criação e gestão de formulários.
              </Typography>
              <Typography>
                Quando o login com Google é usado, a aplicação solicita dados básicos do perfil, como nome e e-mail. Essas informações são armazenadas para associar formulários e sessões ao usuário autenticado.
              </Typography>
              <Typography>
                As respostas enviadas em formulários públicos são armazenadas para consulta pelo administrador responsável pelo formulário.
              </Typography>
              <Typography>
                Nenhum dado é vendido ou compartilhado com terceiros para fins de marketing.
              </Typography>
            </Stack>
          ) : (
            <Stack spacing={2}>
              <Typography>
                O Form Builder é uma aplicação de criação e publicação de formulários. Usuários autenticados são responsáveis pelo conteúdo dos formulários que criam.
              </Typography>
              <Typography>
                Visitantes podem responder formulários publicados por meio de links públicos. O envio de uma resposta indica concordância com o armazenamento das informações preenchidas.
              </Typography>
              <Typography>
                A aplicação é fornecida para fins de avaliação técnica e pode conter limitações documentadas no projeto.
              </Typography>
            </Stack>
          )}
          <Button component={RouterLink} to="/login" variant="outlined">
            Voltar
          </Button>
        </Stack>
      </Paper>
    </Container>
  );
}

function PublicForm({ form }: { form: PublicForm }) {
  const [answers, setAnswers] = useState<Answers>({});
  const fields = useMemo(() => form.fields ?? [], [form.fields]);
  const submitMutation = useMutation({
    mutationFn: async () =>
      unwrap(
        await submitPublicFormResponse({
          path: { slug: form.slug },
          body: { answers },
        }),
      ),
  });

  if (submitMutation.isSuccess) {
    return (
      <Container maxWidth="sm" sx={{ py: 8 }}>
        <Paper elevation={0} className="public-panel">
          <Stack spacing={2} alignItems="flex-start">
            <CheckCircleIcon color="success" fontSize="large" />
            <Typography variant="h4">Resposta enviada</Typography>
            <Typography color="text.secondary">Obrigado por preencher o formulário.</Typography>
          </Stack>
        </Paper>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm" sx={{ py: 8 }}>
      <Paper elevation={0} className="public-panel">
        <Box
          component="form"
          onSubmit={(event) => {
            event.preventDefault();
            submitMutation.mutate();
          }}
        >
          <Stack spacing={3}>
            <Box>
              <Typography variant="h4" fontWeight={700}>
                {form.title}
              </Typography>
              {form.description ? <Typography color="text.secondary">{form.description}</Typography> : null}
            </Box>
            {submitMutation.error ? <Alert severity="error">{submitMutation.error.message}</Alert> : null}
            {fields.map((field) => (
              <PublicField key={field.id} field={field} value={answers[field.id]} onChange={(value) => setAnswers((current) => ({ ...current, [field.id]: value }))} />
            ))}
            <Button type="submit" variant="contained" size="large" startIcon={<SendIcon />} disabled={submitMutation.isPending}>
              Enviar resposta
            </Button>
          </Stack>
        </Box>
      </Paper>
    </Container>
  );
}

function PublicField({ field, value, onChange }: { field: FormField; value: unknown; onChange: (value: unknown) => void }) {
  const commonProps = {
    label: field.label,
    required: field.required,
    helperText: field.helpText,
    fullWidth: true,
  };
  if (field.type === 'textarea') {
    return <TextField {...commonProps} value={String(value ?? '')} onChange={(event) => onChange(event.target.value)} multiline minRows={3} />;
  }
  if (field.type === 'email') {
    return <TextField {...commonProps} value={String(value ?? '')} onChange={(event) => onChange(event.target.value)} type="email" />;
  }
  if (field.type === 'number') {
    return (
      <TextField
        {...commonProps}
        value={value ?? ''}
        onChange={(event) => onChange(event.target.value === '' ? undefined : Number(event.target.value))}
        type="number"
        inputProps={{ min: field.min, max: field.max }}
      />
    );
  }
  if (field.type === 'select') {
    return (
      <FormControl fullWidth required={field.required}>
        <InputLabel>{field.label}</InputLabel>
        <Select label={field.label} value={String(value ?? '')} onChange={(event) => onChange(event.target.value)}>
          {(field.options ?? []).map((option) => (
            <MenuItem key={option} value={option}>
              {option}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    );
  }
  if (field.type === 'checkbox') {
    return <FormControlLabel control={<Checkbox checked={Boolean(value)} onChange={(event) => onChange(event.target.checked)} />} label={field.label} />;
  }
  return <TextField {...commonProps} value={String(value ?? '')} onChange={(event) => onChange(event.target.value)} placeholder={field.placeholder} />;
}

function useCurrentUser() {
  return useQuery({
    queryKey: ['me'],
    queryFn: async () => unwrap(await getCurrentUser()),
    retry: false,
  });
}

function CenteredLoader() {
  return (
    <Box sx={{ display: 'grid', minHeight: '50vh', placeItems: 'center' }}>
      <CircularProgress />
    </Box>
  );
}

function normalizeFieldID(value: string) {
  return value
    .toLowerCase()
    .replace(/[^a-z0-9_]+/g, '_')
    .replace(/^_+|_+$/g, '');
}

function formatAnswer(value: unknown) {
  if (typeof value === 'boolean') {
    return value ? 'Sim' : 'Não';
  }
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return String(value);
}

export default App;
