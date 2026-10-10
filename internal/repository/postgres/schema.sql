-- 1. EXTENSIONS & UUIDv7 GENERATOR
DO $$
BEGIN
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";
EXCEPTION WHEN OTHERS THEN
    NULL;
END $$;

-- Pure PL/pgSQL RFC 9562 UUIDv7 generator
CREATE OR REPLACE FUNCTION uuid_generate_v7() RETURNS uuid AS $$
DECLARE
    unix_time_ms bytea;
    uuid_bytes bytea;
BEGIN
    -- 48 bits (6 bytes) of unix timestamp in milliseconds
    unix_time_ms = substring(int8send(floor(extract(epoch from clock_timestamp()) * 1000)::bigint) from 3 for 6);
    -- 10 random bytes
    uuid_bytes = unix_time_ms || gen_random_bytes(10);
    -- Set version 7 (0111 in high 4 bits of byte 7, index 6 in 0-based)
    uuid_bytes = set_byte(uuid_bytes, 6, (get_byte(uuid_bytes, 6) & 15) | 112);
    -- Set variant 10xx in high 2 bits of byte 9, index 8 in 0-based
    uuid_bytes = set_byte(uuid_bytes, 8, (get_byte(uuid_bytes, 8) & 63) | 128);
    RETURN encode(uuid_bytes, 'hex')::uuid;
END;
$$ LANGUAGE plpgsql VOLATILE;

-- 2. SCHEMA DEFINITION


-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    phone TEXT,
    avatar TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Agents Table
CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    role TEXT,
    city TEXT,
    areas_served JSONB DEFAULT '[]'::jsonb,
    specialization TEXT,
    experience TEXT,
    languages JSONB DEFAULT '[]'::jsonb,
    rera_id TEXT,
    photo TEXT,
    phone TEXT,
    whatsapp TEXT,
    email TEXT,
    active_listings INT DEFAULT 0,
    bio TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Locations Table
CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name TEXT NOT NULL UNIQUE,
    tagline TEXT,
    avg_price TEXT,
    rental_yield TEXT,
    image TEXT,
    key_localities JSONB DEFAULT '[]'::jsonb,
    description TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Projects Table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name TEXT NOT NULL,
    developer TEXT,
    location TEXT,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    lead_agent_id UUID REFERENCES agents(id) ON DELETE SET NULL,
    starting_price TEXT,
    configuration TEXT,
    possession_date TEXT,
    status TEXT,
    badge TEXT,
    total_units TEXT,
    land_area TEXT,
    rera_number TEXT,
    image TEXT,
    overview TEXT,
    highlights JSONB DEFAULT '[]'::jsonb,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Properties Table
CREATE TABLE IF NOT EXISTS properties (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    type TEXT NOT NULL,
    sub_type TEXT,
    transaction TEXT NOT NULL,
    price BIGINT NOT NULL,
    price_display TEXT,
    price_per_sq_ft TEXT,
    deposit TEXT,
    bedrooms INT DEFAULT 0,
    bathrooms INT DEFAULT 0,
    balconies INT DEFAULT 0,
    area INT DEFAULT 0,
    carpet_area INT DEFAULT 0,
    area_unit TEXT DEFAULT 'sq.ft.',
    location TEXT NOT NULL,
    city TEXT NOT NULL,
    state TEXT,
    facing TEXT,
    floor TEXT,
    parking TEXT,
    furnishing TEXT,
    possession TEXT,
    property_age TEXT,
    rera_status TEXT,
    rera_number TEXT,
    featured BOOLEAN DEFAULT FALSE,
    verified BOOLEAN DEFAULT FALSE,
    is_new_launch BOOLEAN DEFAULT FALSE,
    image TEXT,
    gallery JSONB DEFAULT '[]'::jsonb,
    description TEXT,
    amenities JSONB DEFAULT '[]'::jsonb,
    nearby JSONB DEFAULT '[]'::jsonb,
    developer TEXT,
    agent_id UUID REFERENCES agents(id) ON DELETE SET NULL,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Insight Articles Table
CREATE TABLE IF NOT EXISTS insights (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    category TEXT,
    read_time TEXT,
    date TEXT,
    summary TEXT,
    author TEXT,
    author_agent_id UUID REFERENCES agents(id) ON DELETE SET NULL,
    key_takeaway TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Testimonials Table
CREATE TABLE IF NOT EXISTS testimonials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    quote TEXT NOT NULL,
    client TEXT NOT NULL,
    location TEXT,
    type TEXT,
    property TEXT,
    property_id UUID REFERENCES properties(id) ON DELETE SET NULL,
    agent_id UUID REFERENCES agents(id) ON DELETE SET NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Highlights Table
CREATE TABLE IF NOT EXISTS highlights (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    number TEXT,
    title TEXT NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Leads Table
CREATE TABLE IF NOT EXISTS leads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT,
    property_id UUID REFERENCES properties(id) ON DELETE SET NULL,
    property_title TEXT,
    message TEXT,
    type TEXT,
    status TEXT DEFAULT 'pending',
    preferred_date TEXT,
    agent_id UUID REFERENCES agents(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 4. INDEXES
CREATE INDEX IF NOT EXISTS idx_properties_slug ON properties(slug);
CREATE INDEX IF NOT EXISTS idx_properties_city ON properties(city);
CREATE INDEX IF NOT EXISTS idx_properties_type ON properties(type);
CREATE INDEX IF NOT EXISTS idx_properties_transaction ON properties(transaction);
CREATE INDEX IF NOT EXISTS idx_properties_price ON properties(price);
CREATE INDEX IF NOT EXISTS idx_properties_featured ON properties(featured);
CREATE INDEX IF NOT EXISTS idx_properties_agent_id ON properties(agent_id);
CREATE INDEX IF NOT EXISTS idx_properties_project_id ON properties(project_id);
CREATE INDEX IF NOT EXISTS idx_properties_location_id ON properties(location_id);
CREATE INDEX IF NOT EXISTS idx_projects_location_id ON projects(location_id);
CREATE INDEX IF NOT EXISTS idx_projects_lead_agent_id ON projects(lead_agent_id);
CREATE INDEX IF NOT EXISTS idx_agents_city ON agents(city);
CREATE INDEX IF NOT EXISTS idx_agents_user_id ON agents(user_id);
CREATE INDEX IF NOT EXISTS idx_leads_created_at ON leads(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_leads_agent_id ON leads(agent_id);
CREATE INDEX IF NOT EXISTS idx_leads_user_id ON leads(user_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 5. ROW LEVEL SECURITY (RLS)
ALTER TABLE agents ENABLE ROW LEVEL SECURITY;
ALTER TABLE properties ENABLE ROW LEVEL SECURITY;
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE locations ENABLE ROW LEVEL SECURITY;
ALTER TABLE insights ENABLE ROW LEVEL SECURITY;
ALTER TABLE testimonials ENABLE ROW LEVEL SECURITY;
ALTER TABLE highlights ENABLE ROW LEVEL SECURITY;
ALTER TABLE leads ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Public read access policies
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'agents' AND policyname = 'public_read_agents') THEN
        CREATE POLICY public_read_agents ON agents FOR SELECT USING (true);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'properties' AND policyname = 'public_read_properties') THEN
        CREATE POLICY public_read_properties ON properties FOR SELECT USING (true);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'projects' AND policyname = 'public_read_projects') THEN
        CREATE POLICY public_read_projects ON projects FOR SELECT USING (true);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'locations' AND policyname = 'public_read_locations') THEN
        CREATE POLICY public_read_locations ON locations FOR SELECT USING (true);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'insights' AND policyname = 'public_read_insights') THEN
        CREATE POLICY public_read_insights ON insights FOR SELECT USING (true);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'testimonials' AND policyname = 'public_read_testimonials') THEN
        CREATE POLICY public_read_testimonials ON testimonials FOR SELECT USING (true);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'highlights' AND policyname = 'public_read_highlights') THEN
        CREATE POLICY public_read_highlights ON highlights FOR SELECT USING (true);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'leads' AND policyname = 'public_insert_leads') THEN
        CREATE POLICY public_insert_leads ON leads FOR INSERT WITH CHECK (true);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'leads' AND policyname = 'admin_read_leads') THEN
        CREATE POLICY admin_read_leads ON leads FOR SELECT USING (true);
    END IF;
END $$;

-- 6. STORAGE BUCKET CONFIGURATION
DO $$
BEGIN
    BEGIN
        IF EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = 'storage') THEN
            INSERT INTO storage.buckets (id, name, public, file_size_limit, allowed_mime_types)
            VALUES (
                'test-bkt',
                'test-bkt',
                true,
                10485760, -- 10MB limit
                ARRAY['image/jpeg', 'image/png', 'image/webp', 'image/svg+xml']
            )
            ON CONFLICT (id) DO UPDATE SET
                public = true,
                allowed_mime_types = ARRAY['image/jpeg', 'image/png', 'image/webp', 'image/svg+xml'];
        END IF;
    EXCEPTION WHEN OTHERS THEN
        NULL;
    END;
END $$;


-- 7. SCHEMA MIGRATIONS
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
