// Common API types
export type ApiUser = {
  id: string;
  email: string;
  username?: string;
  role: string;
  cpuLimit?: number;
  memoryMbLimit?: number;
  diskMbLimit?: number;
  backupLimit?: number;
  databaseLimit?: number;
  allocationLimit?: number;
  subuserLimit?: number;
  scheduleLimit?: number;
  serverLimit?: number;
  createdAt: string;
  updatedAt: string;
  sessionVersion?: number;
  disabled?: boolean;
  useTotp?: boolean;
  totpSecret?: string;
  image_url?: string;
};

export type ApiServer = {
  id: string;
  name: string;
  description?: string;
  owner?: string;
  ownerId?: string;
  ownerEmail?: string;
  template?: string;
  sftpHost?: string;
  sftpPort?: number;
  permissions?: string[];
  status: string;
  nodeId?: string;
  node?: string;
  allocationId?: string;
  primaryAllocationId?: string;
  allocation?: string;
  cpuLimit?: number;
  cpuShares?: number;
  memoryMb?: number;
  diskMb?: number;
  databaseLimit?: number;
  backupLimit?: number;
  allocationLimit?: number;
  ioWeight?: number;
  suspended?: boolean;
  transferring?: boolean;
  installing?: boolean;
  installationState?: string;
  transferTargetNodeId?: string;
  transferState?: string;
  transferError?: string;
  uuid?: string;
  image?: string;
  createdAt?: string;
  memory?: string | null;
  cpu?: string | null;
  uptime?: string | null;
  featureLimits?: {
    databases: number;
    allocations: number;
    backups: number;
  };
  dockerImage?: string;
  startupCommand?: string;
  environment?: Record<string, string>;
  relationship?: 'owner' | 'subuser' | 'admin';
  configSyncPending?: boolean;
  configSyncError?: string;
};

export type ApiNodeActualState = "online" | "offline" | "degraded" | "unknown";

export type ApiNode = {
  id: string;
  uuid?: string;
  name: string;
  region: string;
  regionId?: string;
  locationId?: string;
  status: string;
  desiredState?: string;
  // `actualState` is the canonical operational state. `status` is legacy
  // configuration data and must not be used to infer node connectivity.
  actualState?: ApiNodeActualState;
  heartbeatState?: string;
  heartbeatStateChangedAt?: string;
  heartbeatRecoveryCount?: number;
  maintenanceMode?: boolean;
  draining?: boolean;
  behindProxy?: boolean;
  placementEligible?: boolean;
  baseUrl?: string;
  fqdn?: string;
  scheme?: string;
  description?: string;
  public?: boolean;
  isPublic?: boolean;
  tokenId?: string;
  daemonBase?: string;
  daemonListen?: number;
  daemonSftp?: number;
  lastSeenAt?: string;
  memoryMb?: number;
  diskMb?: number;
  uploadSizeMb?: number;
  cpu?: number;
  memory?: number;
  disk?: number;
  servers?: number;
  version?: string;
  os?: string;
  architecture?: string;
  cpuThreads?: number;
  dockerStatus?: string;
  nodeMemoryMb?: number;
  nodeDiskMb?: number;
  heartbeatError?: string;
};

export type ApiNodeHealth = {
  cpu: string;
  memory: string;
  disk: string;
  network: string;
  runtime: string;
};

export type ApiNodeHealthScore = {
  cpu: number;
  memory: number;
  disk: number;
  heartbeat: number;
  status: number;
  total: number;
};

export type ApiNodeCapacity = {
  nodeId: string;
  regionId?: string;
  allocated_cpu: number;
  available_cpu: number;
  allocated_memory: number;
  available_memory: number;
  allocated_disk: number;
  available_disk: number;
  server_count: number;
  updated_at: string;
};

export type ApiNodeLifecycle = {
  node: ApiNode;
  health: ApiNodeHealth;
  healthScore: ApiNodeHealthScore;
  capacity: ApiNodeCapacity;
  draining: boolean;
  maintenance: boolean;
  placementEligible: boolean;
  placementBlockedReason?: string;
};

// Minimal node representation exposed by the allocation inventory endpoint.
export type ApiAllocationNode = {
  id: string;
  name: string;
};

export type ApiAllocation = {
  id: string;
  node: string;
  nodeId?: string;
  ip: string;
  port: number;
  alias?: string;
  notes?: string;
  server?: string;
  serverId?: string;
  isPrimary?: boolean;
  primary?: boolean;
};

export type ApiDatabase = {
  id: string;
  serverId: string;
  name: string;
  username: string;
  remote?: string;
  password?: string;
  maxConnections?: number;
  hostId?: string;
  host?: string;
  port?: number;
  database?: string;
  engine?: string;
  provisioningState?: "pending" | "ready" | "failed" | string;
  provisioningError?: string;
  createdAt: string;
  updatedAt: string;
};

export type ApiDatabaseOrphanRemediation = {
  id: string;
  serverDatabaseId: string;
  serverId: string;
  databaseHostId: string;
  engine: string;
  host: string;
  port: number;
  database: string;
  username: string;
  remote: string;
  reason: string;
  status: "pending" | "resolved";
  createdAt: string;
  resolvedAt?: string;
};

export type ApiServerOrphanRemediation = {
  id: string;
  serverId: string;
  nodeUrl: string;
  daemonError: string;
  status: "pending" | "resolved";
  createdAt: string;
  resolvedAt?: string;
};

export type ApiOrphanRemediations = {
  serverRemediations: ApiServerOrphanRemediation[];
  databaseRemediations: ApiDatabaseOrphanRemediation[];
};

export type ApiBackup = {
  id: string;
  serverId: string;
  name: string;
  uuid?: string;
  successful: boolean;
  locked?: boolean;
  isLocked?: boolean;
  status?: string;
  completedAt?: string;
  createdAt: string;
  size?: number;
  checksum?: string;
  statusMessage?: string;
  ignoredFiles?: string[];
};

export type BackupCreateInput = {
  name?: string;
  ignored?: string[];
  is_locked?: boolean;
};

export type ServerCreateInput = {
  name: string;
  description?: string;
  nodeId?: string;
  nestId?: string;
  eggId?: string;
  dockerImage?: string;
  startupCommand?: string;
  environment?: Record<string, string>;
  allocationId?: string;
  ownerId?: string;
  memoryMb?: number;
  diskMb?: number;
  cpuLimit?: number;
  cpuShares?: number;
  regionId?: string;
  templateId?: string;
  limits?: {
    cpu?: number;
    memory?: number;
    disk?: number;
  };
  featureLimits?: {
    databases?: number;
    allocations?: number;
    backups?: number;
  };
};

export type ServerUpdateInput = {
  name?: string;
  description?: string;
  dockerImage?: string;
  startupCommand?: string;
  environment?: Record<string, string>;
  allocationId?: string;
  primaryAllocationId?: string;
  ownerId?: string;
  memoryMb?: number;
  diskMb?: number;
  cpuLimit?: number;
  cpuShares?: number;
  regionId?: string;
  templateId?: string;
  limits?: {
    cpu?: number;
    memory?: number;
    disk?: number;
  };
  featureLimits?: {
    databases?: number;
    allocations?: number;
    backups?: number;
  };
};

export type DatabaseCreateInput = {
  name: string;
  remote?: string;
  hostId?: string;
  username?: string;
  password?: string;
};

export type ScheduleCreateInput = {
  name: string;
  cron?: {
    minute: string;
    hour: string;
    dayOfMonth: string;
    month: string;
    dayOfWeek: string;
  };
  cronMinute?: string;
  cronHour?: string;
  cronDayOfMonth?: string;
  cronMonth?: string;
  cronDayOfWeek?: string;
  isActive?: boolean;
  isProcessing?: boolean;
  onlyWhenOnline?: boolean;
  timezone?: string;
  enabled?: boolean;
};

export type ScheduleUpdateInput = {
  name?: string;
  cron?: {
    minute: string;
    hour: string;
    dayOfMonth: string;
    month: string;
    dayOfWeek: string;
  };
  cronMinute?: string;
  cronHour?: string;
  cronDayOfMonth?: string;
  cronMonth?: string;
  cronDayOfWeek?: string;
  isActive?: boolean;
  isProcessing?: boolean;
  onlyWhenOnline?: boolean;
  timezone?: string;
  enabled?: boolean;
};

export type ScheduleTaskCreateInput = {
  action: string;
  payload?: any;
  continueOnFailure?: boolean;
  timeOffset?: number;
  timeOffsetSeconds?: number;
  sequenceId?: string;
  sequence?: number;
  value?: string;
};

export type ScheduleTaskUpdateInput = {
  action?: string;
  payload?: any;
  continueOnFailure?: boolean;
  timeOffset?: number;
  timeOffsetSeconds?: number;
  sequenceId?: string;
  sequence?: number;
  value?: string;
};

export type ApiSchedule = {
  id: string;
  serverId: string;
  name: string;
  cronMinute: string;
  cronHour: string;
  cronDayOfMonth: string;
  cronMonth: string;
  cronDayOfWeek: string;
  onlyWhenOnline: boolean;
  enabled: boolean;
  timezone?: string;
  createdAt: string;
  updatedAt: string;
  lastRunAt?: string;
  nextRunAt?: string;
  tasks?: ApiScheduleTask[];
};

export type ApiScheduleTask = {
  id: string;
  scheduleId: string;
  action: string;
  payload: Record<string, unknown>;
  continueOnFailure: boolean;
  sequenceOrder: number;
  sequence?: number;
  timeOffset?: number;
  timeOffsetSeconds?: number;
};

export type ApiPublicPanelSettings = {
  companyName: string;
  shortName: string;
  productName: string;
  browserTitle: string;
  footerText: string;
  logoUrl: string;
  faviconUrl: string;
  loginBackgroundUrl: string;
  recaptchaSiteKey?: string;
  recaptchaEnabled?: boolean;
  themePreset: string;
  defaultLocale: string;
};

export type ApiSetupStatus = {
  required: boolean;
  hasAdmin: boolean;
  appVersion: string;
};

export type ApiSetupRequest = {
  email: string;
  password: string;
  name?: string;
};

export type LoginResponse = {
  complete: boolean;
  token?: string;
  user?: ApiUser;
  confirmationToken?: string;
};

export type ApiPanelSettings = {
  companyName: string;
  shortName?: string;
  productName?: string;
  browserTitle?: string;
  footerText?: string;
  logoUrl?: string;
  faviconUrl?: string;
  loginBackgroundUrl?: string;
  themePreset?: string;
  defaultLocale: string;
  require2FA?: "none" | "admin" | "all";
  requireEmailVerification?: boolean;
  passwordComplexity?: string;
  passwordExpirationDays?: number;
  sessionDurationMinutes?: number;
  loginRateLimitEnabled?: boolean;
  loginAttemptThreshold?: number;
  accountLockoutMinutes?: number;
  geoRestrictions?: string;
  apiTokenTtlDays?: number;
  apiRotationDays?: number;
  allowedOrigins?: string;
  trustedNetworks?: string;
  defaultTimezone?: string;
  dateFormat?: string;
  numberFormat?: string;
  currencyFormat?: string;
  defaultDashboard?: string;
  landingPage?: string;
  sidebarLayout?: string;
  compactMode?: boolean;
  advancedMode?: boolean;
  metricsRetentionDays?: number;
  logsRetentionDays?: number;
  auditRetentionDays?: number;
  metricsSamplingRate?: number;
  monitoringPollIntervalSeconds?: number;
  emailAlertsEnabled?: boolean;
  webhookAlertsEnabled?: boolean;
  discordWebhookUrl?: string;
  slackWebhookUrl?: string;
  telegramBotToken?: string;
  placementStrategy?: string;
  antiAffinityRules?: string;
  resourceReservationsEnabled?: boolean;
  nodePrioritization?: string;
  recoveryStrategy?: string;
  failoverThresholdSeconds?: number;
  heartbeatThresholdSeconds?: number;
  reservationDurationMinutes?: number;
  reservationCleanupMinutes?: number;
  capacityBufferPercent?: number;
  backupProvider?: string;
  backupRetentionDays?: number;
  backupLimit?: number;
  backupAutoCleanup?: boolean;
  backupEncryptionEnabled?: boolean;
  backupKeyRotationDays?: number;
};

export type CreateServerDatabaseInput = {
  /** Backend request field; `name` is retained only for legacy callers. */
  database?: string;
  name?: string;
  hostId?: string;
  remote?: string;
  username?: string;
  password?: string;
  maxConnections?: number;
};

export type ApiWSTicket = {
  token: string;
  expiresAt: string;
};

export type ApiNodeConfiguration = {
  id: string;
  nodeId: string;
  config?: string;
  configFormat?: string;
  token?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type ApiDatabaseHostConnectionTestResult = {
  ok: boolean;
  message?: string;
};

export type ApiServerDatabaseDeleteResult = {
  ok: boolean;
  orphanRemediation?: boolean;
};

export type ApiDatabaseHost = {
  id: string;
  name: string;
  host: string;
  port: number;
  username: string;
  database: string;
  engine: string;
  maxConnections?: number;
  maxDatabases?: number;
  tlsMode?: string;
  tlsServerName?: string;
  nodeName?: string;
  databases?: number;
  nodeId?: string;
  linkedNode?: string;
  createdAt: string;
  updatedAt: string;
};

// The `/admin-scopes` endpoint returns a scope-to-description map containing
// only the scopes the current identity may delegate.
export type AdminScopes = Record<string, string>;

export type ApiKey = {
  id: string;
  identifier: string;
  description: string;
  scopes: string[];
  allowedIps?: string[];
  token: string;
  lastUsedAt?: string;
  createdAt: string;
};

export type ApiSSHKey = {
  id?: string;
  name: string;
  publicKey: string;
  fingerprint: string;
  createdAt: string;
};

export type ApiOAuthClient = {
  id: string;
  name: string;
  description?: string;
  scopes: string[];
  allowedScopes: string[];
  scope: string;
  clientId?: string;
  clientSecretHash?: string;
  ownerId?: string;
  createdAt: string;
  updatedAt: string;
};

export type ApiOAuthClientCreation = {
  id: string;
  name: string;
  description?: string;
  scopes: string[];
  clientId?: string;
  clientSecret?: string;
  client?: { id: string; clientId: string; name: string; scopes: string[] };
  createdAt: string;
};

export type ApiPlugin = {
  id: string;
  name: string;
  description?: string;
  version?: string;
  author?: string;
  enabled: boolean;
  kind?: string;
  installedAt: string;
};

export type ApiWebhook = {
  id: string;
  name: string;
  description?: string;
  url: string;
  events: string[];
  enabled: boolean;
  secret?: string;
  webhookType?: string;
  discordUsername?: string;
  discordAvatarUrl?: string;
  discordContent?: string;
  createdAt: string;
  updatedAt: string;
};

export type ApiWebhookDelivery = {
  id: string;
  webhookId: string;
  event: string;
  eventName?: string;
  status: string;
  state?: string;
  statusCode?: number;
  responseStatus?: number;
  responseBody?: string;
  responseBodyExcerpt?: string;
  attempt: number;
  success: boolean;
  lastError?: string;
  createdAt: string;
};

export type ApiMigrationStatus =
  | "pending"
  | "planned"
  | "preparing"
  | "transferring"
  | "restoring"
  | "in_progress"
  | "completed"
  | "failed"
  | "cancelled";

export type ApiMigrationHistory = {
  id: string;
  migrationId: string;
  serverId: string;
  sourceNodeId: string;
  targetNodeId: string;
  status: ApiMigrationStatus;
  fromStatus: string;
  toStatus: string;
  reason: string;
  startedAt?: string;
  completedAt?: string;
  error?: string;
  createdAt: string;
};

export type ApiMigration = {
  id: string;
  serverId: string;
  sourceNodeId: string;
  targetNodeId: string;
  status: ApiMigrationStatus;
  initiatedBy?: string;
  priority?: number;
  transferMethod?: string;
  transferPhase?: string;
  idempotencyKey?: string;
  archiveSize?: number;
  archiveChecksum?: string;
  cleanupPending?: boolean;
  history?: ApiMigrationHistory[];
  progress?: number;
  error?: string;
  failureReason?: string;
  startedAt?: string;
  completedAt?: string;
  createdAt: string;
  updatedAt: string;
};

export type CreateMigrationInput = {
  serverId: string;
  sourceNodeId?: string;
  targetNodeId?: string;
};

export type ApiPanelMailSettings = {
  enabled?: boolean;
  mailFrom?: string;
  mailFromAddress?: string;
  mailFromName?: string;
  host?: string;
  smtpHost?: string;
  port?: number;
  smtpPort?: number;
  encryption?: string;
  smtpEncryption?: string;
  username?: string;
  smtpUsername?: string;
  password?: string;
  smtpPassword?: string;
};

export type ApiPanelAdvancedSettings = {
  reCAPTCHAEnabled?: boolean;
  reCAPTCHASiteKey?: string;
  reCAPTCHASecretKey?: string;
  recaptchaEnabled?: boolean;
  recaptchaWebsiteKey?: string;
  recaptchaSecretKey?: string;
  analyticsCode?: string;
  maintenanceMode?: boolean;
  maintenanceMessage?: string;
  guzzleConnectTimeout?: number;
  guzzleRequestTimeout?: number;
  autoAllocEnabled?: boolean;
  autoAllocStartPort?: number;
  autoAllocEndPort?: number;
};

export type ApiMount = {
  id: string;
  name: string;
  description?: string;
  source: string;
  target: string;
  readOnly: boolean;
  uuid?: string;
  userMountable?: boolean;
  templateIds?: string[];
  nodeIds?: string[];
  serverIds?: string[];
  eggs?: string[];
  nodes?: string[];
  servers?: string[];
  createdAt: string;
  updatedAt: string;
};

export type ApiNest = {
  id: string;
  name: string;
  description?: string;
  eggs?: number;
  eggCount?: number;
  createdAt: string;
  updatedAt: string;
};

export type ApiEgg = {
  id: string;
  nestId: string;
  name: string;
  description?: string;
  dockerImage?: string;
  dockerImages?: string[];
  startupCommand?: string;
  startup?: string;
  config?: any;
  configFiles?: string;
  environment?: Record<string, string>;
  variables?: ApiStartupVariable[];
  installScript?: string;
  installContainer?: string;
  installEntrypoint?: string;
  image?: string;
  defaultMemoryMb?: number;
  nestName?: string;
  createdAt: string;
  updatedAt: string;
};

export type ApiRole = {
  id: string;
  key: string;
  name: string;
  isAdmin: boolean;
  permissions?: string[];
  createdAt?: string;
  updatedAt?: string;
};

export type ApiRegion = {
  id: string;
  uuid: string;
  name: string;
  slug: string;
  description: string;
  enabled: boolean;
  nodeCount: number;
  createdAt: string;
  updatedAt: string;
};

export type ApiLocation = {
  id: string;
  short: string;
  long: string;
  nodeCount: number;
  serverCount: number;
  createdAt: string;
};

export type ApiAdminAuditEvent = {
  id: string;
  userId: string;
  userEmail?: string;
  action: string;
  resource: string;
  resourceId?: string;
  metadata?: Record<string, unknown>;
  ip?: string;
  userAgent?: string;
  actorEmail?: string;
  targetType?: string;
  targetId?: string;
  createdAt: string;
};

export type ApiStats = {
  cpuPercent: number;
  memoryBytes: number;
  memoryLimit: number;
  diskBytes: number;
  diskLimit: number;
  networkRxBytes: number;
  networkTxBytes: number;
  uptime: number;
  uptimeMs?: number;
  state?: string;
};

export type TwoFactorSetup = {
  enabled: boolean;
  secret?: string;
  qrCodeUrl?: string;
  qrCode?: string;
  image_url?: string;
  tokens?: string[];
};

export type CreateNodeInput = {
  name: string;
  region: string;
  regionId?: string;
  locationId: string;
  fqdn: string;
  scheme?: string;
  behindProxy?: boolean;
  public?: boolean;
  daemonBase?: string;
  daemonListen?: number;
  daemonSftp?: number;
  description?: string;
  memoryMb?: number;
  diskMb?: number;
  uploadSizeMb?: number;
  placementEligible?: boolean;
  baseUrl?: string;
  status?: string;
  desiredState?: string;
  maintenanceMode?: boolean;
  draining?: boolean;
};

export type UpdateNodeInput = {
  name?: string;
  region?: string;
  regionId?: string;
  locationId?: string;
  fqdn?: string;
  scheme?: string;
  behindProxy?: boolean;
  public?: boolean;
  daemonBase?: string;
  daemonListen?: number;
  daemonSftp?: number;
  description?: string;
  memoryMb?: number;
  diskMb?: number;
  uploadSizeMb?: number;
  placementEligible?: boolean;
  baseUrl?: string;
  status?: string;
  desiredState?: string;
  maintenanceMode?: boolean;
  draining?: boolean;
};

export type CreateAllocationInput = {
  nodeId: string;
  ip: string;
  ports: string;
  alias?: string;
  notes?: string;
};

export type UpdateAllocationInput = {
  alias?: string;
  notes?: string;
};

export type CreateDatabaseHostInput = {
  name: string;
  host: string;
  port: number;
  username: string;
  password: string;
  engine?: string;
  tlsMode?: string;
  tlsServerName?: string;
  tlsCa?: string;
  nodeId?: string;
  maxDatabases?: number;
};

export type UpdateDatabaseHostInput = {
  name?: string;
  host?: string;
  port?: number;
  username?: string;
  password?: string;
  engine?: string;
  tlsMode?: string;
  tlsServerName?: string;
  tlsCa?: string;
  nodeId?: string;
  maxDatabases?: number;
};

export type CreateMountInput = {
  name: string;
  source: string;
  target: string;
  readOnly: boolean;
  description?: string;
  userMountable?: boolean;
  nodeIds?: string[];
  templateIds?: string[];
};

export type AssignMountInput = {
  mountId: string;
  readOnly?: boolean;
};

export type ApiMountAssignmentResponse = {
  ok: boolean;
  runtimeSynchronized: boolean;
};

export type RenameFileInput = {
  from: string;
  to: string;
};

export type PatchScheduleTaskInput = {
  action?: string;
  payload?: any;
  continueOnFailure?: boolean;
  timeOffset?: number;
  timeOffsetSeconds?: number;
  sequence?: number;
  value?: string;
};

export type CreateEggInput = {
  name: string;
  nestId: string;
  description?: string;
  dockerImage?: string;
  dockerImages?: string[];
  startupCommand?: string;
  startup?: string;
  config?: any;
  environment?: Record<string, string>;
  configFiles?: string;
  installScript?: string;
  installContainer?: string;
  installEntrypoint?: string;
};

export type UpdateEggInput = {
  name?: string;
  description?: string;
  dockerImage?: string;
  dockerImages?: string[];
  startupCommand?: string;
  startup?: string;
  config?: any;
  environment?: Record<string, string>;
  configFiles?: string;
  installScript?: string;
  installContainer?: string;
  installEntrypoint?: string;
};

export type ApiActivityLog = {
  id: string;
  userId?: string;
  // Legacy fallback fields; canonical /activity/events responses use event and subject fields.
  action?: string;
  resource?: string;
  description?: string;
  event: string;
  timestamp: string;
  resourceId?: string;
  metadata?: Record<string, unknown>;
  ip?: string;
  userAgent?: string;
  actorEmail?: string;
  subjectType?: string;
  subjectId?: string;
  properties?: any;
  level?: string;
  source?: string;
  createdAt?: string;
};

export type ApiAuditEvent = {
  id: string;
  userId?: string;
  actorEmail?: string;
  action: string;
  targetType?: string;
  targetId?: string;
  metadata?: Record<string, unknown>;
  ip?: string;
  createdAt: string;
};

export type ApiFileEntry = {
  name: string;
  path: string;
  directory: boolean;
  size?: number;
  mode?: string;
  mime?: string;
  symlink?: boolean;
  modifiedAt?: string;
  createdAt?: string;
};

export type ApiServerSubuser = {
  id: string;
  userId?: string;
  email: string;
  permissions: string[];
  createdAt?: string;
  updatedAt?: string;
};

export type ApiStartupVariable = {
  id: string;
  name: string;
  description?: string;
  envVariable: string;
  env_variable?: string;
  defaultValue: string;
  default_value?: string;
  serverValue: string;
  server_value?: string;
  rules: string;
  is_editable?: boolean;
};

export type CrashEvent = {
  id: string;
  server_id: string;
  node_id: string;
  exit_code: number;
  oom_killed: boolean;
  clean_exit: boolean;
  auto_restarted: boolean;
  crash_count: number;
  node_state: Record<string, unknown> | null;
  created_at: string;
};

export type ApiHealthCheck = {
  name: string;
  status: "ok" | "warning" | "failed";
  label: string;
  notificationMessage: string;
  critical: boolean;
  latencyMs?: number;
  details?: Record<string, unknown>;
  lastChecked?: string;
  lastSuccess?: string;
  lastFailure?: string;
  consecutiveFailures?: number;
};

// Returned by the diagnostic endpoint (`GET /health`). Unlike readiness and
// liveness, this endpoint always returns a report; callers must inspect status.
export type ApiHealthReport = {
  status: "ok" | "warning" | "failed";
  ok: boolean;
  service: string;
  version?: string;
  uptime?: string;
  checks: ApiHealthCheck[];
  checkedAt: string;
};

export type ApiTemplate = {
  id: string;
  name: string;
  description?: string;
  eggId: string;
  nestId: string;
  dockerImage?: string;
  startupCommand?: string;
  environment?: Record<string, string>;
  createdAt: string;
  updatedAt: string;
};

export type ApiUserSearchResult = {
  users: ApiUser[];
  total: number;
  page: number;
  perPage: number;
};

export type ApiUserSession = {
  id: string;
  uuid?: string;
  userId?: string;
  ip?: string;
  ipAddress: string;
  userAgent?: string;
  lastActivity?: string;
  createdAt?: string;
  current?: boolean;
  isRevoked?: boolean;
  revokedAt?: string;
  revokeReason?: string;
  expiresAt?: string;
};

export type ApiEvacuationItem = {
  id: string;
  planId: string;
  serverId: string;
  sourceNodeId: string;
  targetNodeId?: string;
  eligible: boolean;
  reason: string;
  migrationId?: string;
  status: string;
  error?: string;
};

export type ApiEvacuationPlan = {
  id: string;
  nodeId: string;
  status: "pending" | "running" | "completed" | "cancelled" | "failed";
  items: ApiEvacuationItem[];
  createdAt: string;
  updatedAt: string;
};

export type ApiEvacuationResult = {
  plan: ApiEvacuationPlan;
  items: ApiEvacuationItem[];
  preview: boolean;
};

export type ApiRecoveryItem = {
  id: string;
  planId: string;
  serverId: string;
  sourceNodeId: string;
  targetNodeId?: string;
  reservationId?: string;
  migrationId?: string;
  sourceBackupName?: string;
  sourceBackupChecksum?: string;
  sourceBackupSize?: number;
  status: "pending" | "planned" | "executing" | "completed" | "restored" | "cancelled" | "failed" | "skipped";
  reason?: string;
  createdAt: string;
  updatedAt: string;
};

export type ApiReservation = {
  id: string;
  nodeId: string;
  status: string;
  serverId?: string;
  migrationId?: string;
  reservationType?: string;
  reservedBy?: string;
  cpu?: number;
  memory?: number;
  disk?: number;
  createdAt?: string;
  updatedAt?: string;
};

export type ApiRecoveryPlan = {
  id: string;
  nodeId: string;
  status: "pending" | "planning" | "planned" | "executing" | "completed" | "restored" | "cancelled" | "failed";
  reason: string;
  items: ApiRecoveryItem[];
  createdAt: string;
  updatedAt: string;
};

export type CreateRecoveryPlanInput = {
  nodeId: string;
  reason: string;
};

/** Read-only compatibility view for the retired server transfer API. */
export type ApiLegacyTransferStatus = {
  state?: string;
  transferring: boolean;
  targetNodeId?: string;
  error?: string;
};

export type SocialProvider = {
  id: string;
  name: string;
  displayName: string;
  enabled: boolean;
  clientId: string;
  issuerUrl?: string;
  hasClientSecret: boolean;
  scopes: string[];
  buttonStyle: string;
  iconClass: string;
  createdAt: string;
  updatedAt: string;
};
