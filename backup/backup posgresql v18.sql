--
-- PostgreSQL database dump
--

\restrict WzY8voleUmHuD7wmDVSih8kfCIjUbeh3UhFHvfbfN0XZDusWkfPtGy31WBv69X4

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.items (
    id uuid NOT NULL,
    name character varying(150) NOT NULL,
    category_id integer NOT NULL,
    rack_id integer NOT NULL,
    stock integer DEFAULT 1 NOT NULL,
    min_stock integer DEFAULT 1 NOT NULL,
    price numeric(12,2) NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT items_min_stock_check CHECK ((min_stock >= 0)),
    CONSTRAINT items_price_check CHECK ((price > (0)::numeric)),
    CONSTRAINT items_stock_check CHECK ((stock >= 0))
);


ALTER TABLE public.items OWNER TO postgres;

--
-- Name: racks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.racks (
    id integer NOT NULL,
    warehouse_id integer NOT NULL,
    code character varying(50) NOT NULL,
    description text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.racks OWNER TO postgres;

--
-- Name: racks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.racks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.racks_id_seq OWNER TO postgres;

--
-- Name: racks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.racks_id_seq OWNED BY public.racks.id;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.roles_id_seq OWNER TO postgres;

--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- Name: sale_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sale_items (
    id integer NOT NULL,
    sale_id uuid NOT NULL,
    item_id uuid NOT NULL,
    quantity integer NOT NULL,
    price numeric(12,2) NOT NULL,
    subtotal numeric(14,2) NOT NULL,
    CONSTRAINT sale_items_price_check CHECK ((price > (0)::numeric)),
    CONSTRAINT sale_items_quantity_check CHECK ((quantity > 0)),
    CONSTRAINT sale_items_subtotal_check CHECK ((subtotal > (0)::numeric))
);


ALTER TABLE public.sale_items OWNER TO postgres;

--
-- Name: sale_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sale_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sale_items_id_seq OWNER TO postgres;

--
-- Name: sale_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sale_items_id_seq OWNED BY public.sale_items.id;


--
-- Name: sales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sales (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    total_amount numeric(14,2) NOT NULL,
    status character varying(20) DEFAULT 'COMPLETED'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT sales_total_amount_check CHECK ((total_amount > (0)::numeric))
);


ALTER TABLE public.sales OWNER TO postgres;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    revoked_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- Name: stock_adjustments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_adjustments (
    id uuid NOT NULL,
    item_id uuid NOT NULL,
    user_id uuid NOT NULL,
    change integer NOT NULL,
    reason text,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT stock_adjustments_change_check CHECK ((change <> 0))
);


ALTER TABLE public.stock_adjustments OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    password_hash text NOT NULL,
    role_id integer NOT NULL,
    is_active boolean DEFAULT true,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: warehouses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.warehouses (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    location text,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.warehouses OWNER TO postgres;

--
-- Name: warehouses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.warehouses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.warehouses_id_seq OWNER TO postgres;

--
-- Name: warehouses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.warehouses_id_seq OWNED BY public.warehouses.id;


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: racks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks ALTER COLUMN id SET DEFAULT nextval('public.racks_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: sale_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items ALTER COLUMN id SET DEFAULT nextval('public.sale_items_id_seq'::regclass);


--
-- Name: warehouses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses ALTER COLUMN id SET DEFAULT nextval('public.warehouses_id_seq'::regclass);


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, name, description, created_at, updated_at, deleted_at) FROM stdin;
4	PC AIO	Komputer All in One	2026-01-02 06:36:49.711915+07	2026-01-02 06:36:49.711915+07	\N
5	Peripheral	Perangkat Peripheral Komputer	2026-01-02 06:38:20.856482+07	2026-01-02 06:38:20.856482+07	\N
6	Networking	Perangkat Jaringan Komputer	2026-01-02 06:39:26.912557+07	2026-01-02 06:39:26.912557+07	\N
1	Laptop	Laptop dan Notebook	2026-01-02 06:33:43.710941+07	2026-01-02 07:48:18.140417+07	\N
2	Tablet	Tablet dan Aksesoris	2026-01-02 06:35:03.011217+07	2026-01-02 08:03:27.395051+07	\N
7	Motherboard	Motherboard komputer	2026-01-02 12:38:14.8648+07	2026-01-02 12:38:14.8648+07	\N
9	VGA	Komponen Grapics Card Komputer	2026-01-02 12:40:02.770942+07	2026-01-02 12:40:02.770942+07	\N
3	CPU	Komponen Prosessor PC	2026-01-02 06:36:26.29005+07	2026-01-02 06:36:26.29005+07	\N
8	Memori	Komponen Memori Komputer (RAM, Memori Internal)	2026-01-02 12:39:02.302618+07	2026-01-02 12:39:02.302618+07	\N
\.


--
-- Data for Name: items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.items (id, name, category_id, rack_id, stock, min_stock, price, created_at, updated_at, deleted_at) FROM stdin;
7dfb6233-b278-45ee-beff-f5196342cfeb	Intel Core i3 10100F / 10105F Box Comet Lake	3	2	100	10	1350000.00	2026-01-02 12:35:11.802987+07	2026-01-02 12:35:11.802987+07	\N
6b362e3f-1c3b-489c-9aea-4114274e3334	Intel Core i5 12400F 2.5Ghz Up To 4.4Ghz Box Alder Lake	3	2	500	100	2470000.00	2026-01-02 12:43:56.822699+07	2026-01-02 12:43:56.822699+07	\N
59c5243b-607d-4c2e-909a-40f7c5d5bdf7	Intel Core Ultra 9 285K Box Arrow Lake LGA 1851	3	2	20	5	10000000.00	2026-01-02 12:46:26.211468+07	2026-01-02 12:46:26.211468+07	\N
0a1fd854-8e9f-4721-a47d-2319c82f9bf6	AMD AM5 Ryzen 9 9950X3D 4.3Ghz Up To 5.7Ghz	3	2	10	2	12770000.00	2026-01-02 12:49:07.385606+07	2026-01-02 12:49:07.385606+07	\N
5ee98d10-f86b-4032-b929-44136c06578d	MSI Pro B760M-A Wifi DDR5	7	2	100	20	2400000.00	2026-01-02 12:54:41.729715+07	2026-01-02 12:54:41.729715+07	\N
fa6cff85-86d2-42df-ad27-d7e38cbc472e	MSI MAG B760M Gaming Plus Wifi DDR5	7	2	100	20	2700000.00	2026-01-02 12:55:10.617469+07	2026-01-02 12:55:10.617469+07	\N
a4ffa643-2179-4108-b27d-e32803b9739b	MSI Pro B650M-A Wifi	7	2	100	20	2800000.00	2026-01-02 12:55:43.662991+07	2026-01-02 12:55:43.662991+07	\N
488bd09a-eb41-410a-8854-5b5b0bc23c3e	Asrock B550M Steel Legend	7	2	100	20	2200000.00	2026-01-02 12:56:14.25143+07	2026-01-02 12:56:14.25143+07	\N
83c6b4ba-0090-4074-9a19-c880072b62de	Asrock B760M Steel Legend Wifi DDR5	7	2	100	20	2955000.00	2026-01-02 12:56:49.908335+07	2026-01-02 12:56:49.908335+07	\N
2f56643d-52a9-49b8-8062-4cd53f78adb6	Laptop ASUS Core i3 RAM 8GB	1	6	100	30	600000.00	2026-01-02 13:05:25.537046+07	2026-01-02 13:05:25.537046+07	\N
98005180-dba9-4301-992b-628244a8ab52	Laptop ASUS Core i3 RAM 8GB	1	1	700	100	600000.00	2026-01-02 13:05:41.866947+07	2026-01-02 13:05:41.866947+07	\N
146dacd9-3f4e-468d-9c19-704d7050a645	AOC Agon 24 AG241QG 165Hz QHD Monitor	5	3	700	100	2200000.00	2026-01-02 13:07:37.417132+07	2026-01-02 13:07:37.417132+07	\N
422e42d9-b90d-4783-812e-58a1edc44032	AOC Agon 24 AG241QG 165Hz QHD Monitor	5	9	500	100	2200000.00	2026-01-02 13:08:20.597568+07	2026-01-02 13:08:20.597568+07	\N
41e8d19b-9fdf-4f76-b789-902b13707ecc	SAMSUNG 27 Odyssey 27AG320-G3 165Hz 1ms LED	5	9	500	100	2450000.00	2026-01-02 13:09:39.621508+07	2026-01-02 13:09:39.621508+07	\N
484f7435-c10f-48f2-82fd-9edaa88c0057	SAMSUNG 27 Odyssey 27AG320-G3 165Hz 1ms LED	5	3	700	100	2450000.00	2026-01-02 13:09:49.652396+07	2026-01-02 13:09:49.652396+07	\N
ad9aa864-3a01-4189-9bea-b6ac2d1101de	Mikrotik RB450 - Router 5 Port 10/100	6	4	1000	200	1350000.00	2026-01-02 13:12:21.428923+07	2026-01-02 13:12:21.428923+07	\N
4e680ecb-444a-42e0-953e-0686f66a7da9	Mikrotik RB450 - Router 5 Port 10/100	6	5	500	200	1350000.00	2026-01-02 13:12:33.654139+07	2026-01-02 13:12:33.654139+07	\N
81080954-7b25-4325-8fcd-546411cd1bfd	Mikrotik RB-750GR3 (hEX) - Router 5 Port 10/100/1000	6	5	1000	200	875000.00	2026-01-02 13:13:32.484335+07	2026-01-02 13:13:32.484335+07	\N
addf4edb-440f-4dc0-9d3a-6ff052f7a16f	AMD AM4 Ryzen 5 5600g 6 core 12 threads	3	2	100	20	2200000.00	2026-01-02 10:34:39.981777+07	2026-01-02 20:55:24.397347+07	\N
762d72dd-7d8d-44a6-b89e-8527d6c3f10c	Mikrotik RB-750GR3 (hEX) - Router 5 Port 10/100/1000	6	4	1150	200	875000.00	2026-01-02 13:13:46.170473+07	2026-01-03 00:47:19.966818+07	\N
074fbd9d-8c59-43d2-bdf0-2d1ee67c576f	Laptop ASUS Core i5 RAM 8GB	1	6	30	30	7500000.00	2026-01-02 13:05:01.555472+07	2026-01-03 21:37:49.281878+07	\N
19f0144f-a6df-47fd-9225-4aade46b7bff	Laptop ASUS Core i5 RAM 8GB	1	1	450	50	7500000.00	2026-01-02 13:04:33.694775+07	2026-01-03 21:40:53.559314+07	\N
\.


--
-- Data for Name: racks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.racks (id, warehouse_id, code, description, is_active, created_at, updated_at) FROM stdin;
2	1	EL-PC	Rak PC dan Komponen Desktop	t	2026-01-01 00:41:42.598699+07	2026-01-01 00:41:42.598699+07
3	1	EL-ACC	Rak Aksesoris komputer (mouse, keyboard, dll)	t	2026-01-01 00:42:25.583974+07	2026-01-01 00:42:25.583974+07
4	1	EL-NET	Rak Perangkat Jaringan (router, switch, dll)	t	2026-01-01 00:43:01.611779+07	2026-01-01 00:43:01.611779+07
5	2	EL-NET	Rak Perangkat Jaringan (router, switch, dll)	t	2026-01-01 16:59:49.433055+07	2026-01-01 16:59:49.433055+07
7	2	EL-PC	Rak PC dan Komponen Desktop	t	2026-01-01 17:02:05.586002+07	2026-01-01 17:02:05.586002+07
8	3	EL-PC	Rak PC dan Komponen Desktop	t	2026-01-01 17:02:31.344891+07	2026-01-01 17:02:31.344891+07
9	3	EL-ACC	Rak Aksesoris Komputer (mouse, keyboard, dll)	t	2026-01-01 17:03:11.670968+07	2026-01-01 17:03:11.670968+07
1	1	EL-LAP	Rak Laptop, Notebook dan Tablet	t	2026-01-01 00:40:52.169428+07	2026-01-01 22:04:34.018972+07
6	2	EL-LAP	Rak Laptop, Notebook dan Tablet	t	2026-01-01 17:00:31.442562+07	2026-01-01 22:05:32.483583+07
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.roles (id, name, created_at) FROM stdin;
1	super_admin	2025-12-30 13:24:38.034846+07
2	admin	2025-12-30 13:24:38.034846+07
3	staff	2025-12-30 13:24:38.034846+07
\.


--
-- Data for Name: sale_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sale_items (id, sale_id, item_id, quantity, price, subtotal) FROM stdin;
4	0cc5f439-5d77-4d4e-9690-9e091d35d77b	762d72dd-7d8d-44a6-b89e-8527d6c3f10c	50	875000.00	43750000.00
5	709b033b-fda3-4031-98ba-3407bcdb941e	074fbd9d-8c59-43d2-bdf0-2d1ee67c576f	10	7500000.00	75000000.00
7	302fdb1a-3897-4f92-9901-596c0c980ea6	074fbd9d-8c59-43d2-bdf0-2d1ee67c576f	60	7500000.00	450000000.00
8	70feb71c-a674-45d8-9dc7-0f22fef56986	19f0144f-a6df-47fd-9225-4aade46b7bff	50	7500000.00	375000000.00
\.


--
-- Data for Name: sales; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sales (id, user_id, total_amount, status, created_at) FROM stdin;
0cc5f439-5d77-4d4e-9690-9e091d35d77b	d9fac5de-35c8-4c8c-a47e-354b5e12e530	43750000.00	COMPLETED	2026-01-03 00:47:19.966818+07
709b033b-fda3-4031-98ba-3407bcdb941e	24e4fc2c-fda4-48d7-891f-a06dcfeccee6	75000000.00	COMPLETED	2026-01-03 21:01:11.637363+07
302fdb1a-3897-4f92-9901-596c0c980ea6	24e4fc2c-fda4-48d7-891f-a06dcfeccee6	450000000.00	COMPLETED	2026-01-03 21:37:31.319441+07
70feb71c-a674-45d8-9dc7-0f22fef56986	24e4fc2c-fda4-48d7-891f-a06dcfeccee6	375000000.00	COMPLETED	2026-01-03 21:39:45.358087+07
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sessions (id, user_id, expires_at, revoked_at, created_at) FROM stdin;
7bbab15d-245d-42c2-9ce6-b55da3b6e1a6	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 11:31:46.782073+07	2025-12-31 11:18:28.253126+07	2025-12-31 11:16:46.783137+07
ce44b0f7-2f5e-4ac4-9199-7e434cd18e0a	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 11:32:05.680959+07	2025-12-31 11:18:28.253126+07	2025-12-31 11:17:05.682123+07
0058a1f5-3a47-4d34-89cb-8250513185b4	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 11:09:35.352469+07	2025-12-31 11:21:20.720637+07	2025-12-31 10:54:35.353183+07
6ba83c32-30b0-48df-926a-086bf5940886	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 12:13:33.24489+07	\N	2025-12-31 11:58:33.246225+07
4e338e21-a2ee-43e5-b54e-59d1187c4e11	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 12:39:04.50666+07	2025-12-31 12:26:39.388253+07	2025-12-31 12:24:04.50789+07
09d0b5fa-520d-4054-a682-23ba49ff1a82	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 12:41:39.392962+07	\N	2025-12-31 12:26:39.393956+07
36c636dc-ddd5-4579-b597-6b7e553538ac	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 13:21:22.061043+07	2025-12-31 13:11:13.274896+07	2025-12-31 13:06:22.062444+07
ca6f080c-fb29-4746-a7f1-f1fda294ddfb	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 14:55:16.306007+07	2025-12-31 14:40:46.692791+07	2025-12-31 14:40:16.307177+07
3c443d49-3780-4970-b9a0-1fd3cc7764d2	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 14:58:36.986105+07	2025-12-31 14:44:51.34207+07	2025-12-31 14:43:36.987418+07
593bfe4f-06bb-47fa-acb0-17208d4c2040	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 14:59:51.350685+07	2025-12-31 14:46:29.150928+07	2025-12-31 14:44:51.351073+07
5462a3c8-69d5-4b26-9b38-10fde946b4d2	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 15:01:29.157823+07	2025-12-31 14:53:27.274419+07	2025-12-31 14:46:29.158216+07
67edb250-b98e-45e3-b22b-553bd8db8b9b	9c6dc8c8-58bb-47e1-a5e1-32104371b643	2025-12-31 15:08:27.280422+07	\N	2025-12-31 14:53:27.280927+07
d9402f71-c560-43f6-a78e-5ab34c116edf	0f3e9c6a-8f1a-4e2c-9a3b-9a1c7f1c0001	2026-01-01 15:52:06.852481+07	2025-12-31 18:42:30.386115+07	2025-12-31 15:52:06.85309+07
c7b987cc-a487-4769-be9d-bd78a8ce3a73	0f3e9c6a-8f1a-4e2c-9a3b-9a1c7f1c0001	2026-01-01 18:42:30.397765+07	2025-12-31 22:24:31.071479+07	2025-12-31 18:42:30.399731+07
320351a1-326f-42ed-8cff-49e6bb46b47e	0f3e9c6a-8f1a-4e2c-9a3b-9a1c7f1c0001	2026-01-01 22:24:31.078804+07	\N	2025-12-31 22:24:31.083157+07
1aa54321-d0a8-43f9-bdd2-1b8fbf50edaa	d9fac5de-35c8-4c8c-a47e-354b5e12e530	2026-01-01 20:15:16.762692+07	2025-12-31 22:24:39.744042+07	2025-12-31 20:15:16.764018+07
e21ed483-680d-4d4f-9b15-fa02887626e7	24e4fc2c-fda4-48d7-891f-a06dcfeccee6	2026-01-01 20:16:06.046608+07	2025-12-31 22:24:47.305561+07	2025-12-31 20:16:06.049605+07
ff9bee24-6adf-4032-a80c-803a088fde47	24e4fc2c-fda4-48d7-891f-a06dcfeccee6	2026-01-01 22:24:47.305609+07	2025-12-31 22:25:34.032293+07	2025-12-31 22:24:47.306688+07
2760dbd6-73f9-453f-a769-6b5d1d4da2a9	24e4fc2c-fda4-48d7-891f-a06dcfeccee6	2026-01-01 22:25:34.037838+07	\N	2025-12-31 22:25:34.038496+07
c61460d7-957c-4fa2-ba76-713e9fb1e446	d9fac5de-35c8-4c8c-a47e-354b5e12e530	2026-01-01 22:24:39.745614+07	2026-01-01 00:37:27.643593+07	2025-12-31 22:24:39.746266+07
5396cae4-2b24-442d-9f8f-5aa6f09dc8ed	d9fac5de-35c8-4c8c-a47e-354b5e12e530	2026-01-02 00:37:27.687709+07	\N	2026-01-01 00:37:27.689373+07
6144a988-8ef3-4ef1-811f-461e51966a3f	d9fac5de-35c8-4c8c-a47e-354b5e12e530	2026-01-03 06:33:25.497863+07	\N	2026-01-02 06:33:25.498957+07
9d9fea97-9a8f-4201-950d-f7bf435d427f	24e4fc2c-fda4-48d7-891f-a06dcfeccee6	2026-01-04 18:31:47.592024+07	\N	2026-01-03 18:31:47.593566+07
52138aa6-614d-48aa-9a2d-71c0ebc52656	d9fac5de-35c8-4c8c-a47e-354b5e12e530	2026-01-04 12:38:44.67251+07	2026-01-03 18:32:18.315461+07	2026-01-03 12:38:44.674144+07
3cd64351-af07-4df3-a148-a7377d6a1e44	d9fac5de-35c8-4c8c-a47e-354b5e12e530	2026-01-04 18:32:18.318814+07	\N	2026-01-03 18:32:18.319149+07
f9eccee2-8165-4688-89f5-29d4c5db8d1e	0f3e9c6a-8f1a-4e2c-9a3b-9a1c7f1c0001	2026-01-04 18:32:42.857209+07	\N	2026-01-03 18:32:42.857621+07
\.


--
-- Data for Name: stock_adjustments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stock_adjustments (id, item_id, user_id, change, reason, created_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, password_hash, role_id, is_active, deleted_at, created_at, updated_at) FROM stdin;
24e4fc2c-fda4-48d7-891f-a06dcfeccee6	staff-1	endrif@gmail.com	$2a$10$LSypuRx3kCFQ5KpnCJcuLe1nB8nKgDtYkAhmPJ7BD41zRVBjPZ/Ee	3	t	\N	2025-12-30 15:34:20.456561+07	2025-12-30 15:34:20.456561+07
d9fac5de-35c8-4c8c-a47e-354b5e12e530	admin-1	bayu19fr@gmail.com	$2a$10$tYyRCwPKLmPC1Pez2lORN.uXh.vewMExjFdzgUsTKUPH2RjD73xOO	2	t	\N	2025-12-30 15:34:49.151817+07	2025-12-30 15:34:49.151817+07
5e225dc4-a687-497c-9c1f-b31f7d529647	admin-2	silfi@gmail.com	$2a$10$su9SdYHeuOyoqtaDoWZgXescGapad8IYkEgwhH1h1f3DI2G1ELDGy	2	t	\N	2025-12-30 15:35:07.349124+07	2025-12-30 15:35:07.349124+07
9c6dc8c8-58bb-47e1-a5e1-32104371b643	staff-2	rohman@gmail.com	$2a$10$ifLUcC68Xc6hFFMeGkUjxOTiSFPY6xkJyCJfr04kn/w2nQDZ1NGzm	3	t	\N	2025-12-30 15:49:47.819856+07	2025-12-30 15:49:47.819856+07
b7b74b76-6945-4a73-a098-7b94b1ff45c2	staff-3	rohman3@gmail.com	$2a$10$2FKBvivp0/0TjzIBxohuqeG83MPB77vSMXlL1gOT5sTnWUIwyKCIW	3	t	\N	2025-12-31 10:38:22.657726+07	2025-12-31 10:38:22.657726+07
0f3e9c6a-8f1a-4e2c-9a3b-9a1c7f1c0001	Super Administrator	superadmin@ims.local	$2a$10$aK9JDaC1LXF7x00BVWnFa.1peI1jJxBCxEvwM97JFFtY9EJC3uU1y	1	t	\N	2025-12-30 13:51:10.178379+07	2025-12-30 13:51:10.178379+07
\.


--
-- Data for Name: warehouses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.warehouses (id, name, location, is_active, created_at, updated_at) FROM stdin;
1	Gudang Utama	Jakarta	t	2025-12-31 15:54:02.625956+07	2025-12-31 15:54:02.625956+07
2	Gudang Cabang Barat	Bandung	t	2025-12-31 15:55:53.901749+07	2025-12-31 15:55:53.901749+07
4	Gudang Cabang Jawa Tengah	Semarang	f	2025-12-31 18:42:53.139542+07	2025-12-31 18:45:53.818086+07
3	Gudang Cabang Jawa Timur	Kota Surabaya	t	2025-12-31 15:56:03.950699+07	2025-12-31 21:06:52.757418+07
\.


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 9, true);


--
-- Name: racks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.racks_id_seq', 10, true);


--
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_seq', 3, true);


--
-- Name: sale_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sale_items_id_seq', 8, true);


--
-- Name: warehouses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.warehouses_id_seq', 4, true);


--
-- Name: categories categories_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: items items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);


--
-- Name: racks racks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT racks_pkey PRIMARY KEY (id);


--
-- Name: racks racks_warehouse_id_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT racks_warehouse_id_code_key UNIQUE (warehouse_id, code);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: sale_items sale_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_pkey PRIMARY KEY (id);


--
-- Name: sales sales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: stock_adjustments stock_adjustments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_adjustments
    ADD CONSTRAINT stock_adjustments_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: warehouses warehouses_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_name_key UNIQUE (name);


--
-- Name: warehouses warehouses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_pkey PRIMARY KEY (id);


--
-- Name: idx_items_category_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_category_id ON public.items USING btree (category_id);


--
-- Name: idx_items_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_deleted_at ON public.items USING btree (deleted_at);


--
-- Name: idx_items_rack_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_rack_id ON public.items USING btree (rack_id);


--
-- Name: idx_items_stock; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_stock ON public.items USING btree (stock);


--
-- Name: idx_sale_items_item_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sale_items_item_id ON public.sale_items USING btree (item_id);


--
-- Name: idx_sale_items_sale_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sale_items_sale_id ON public.sale_items USING btree (sale_id);


--
-- Name: idx_sales_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_created_at ON public.sales USING btree (created_at);


--
-- Name: idx_sales_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_user_id ON public.sales USING btree (user_id);


--
-- Name: idx_sessions_expires_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sessions_expires_at ON public.sessions USING btree (expires_at);


--
-- Name: idx_sessions_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sessions_user_id ON public.sessions USING btree (user_id);


--
-- Name: idx_stock_adjustments_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_stock_adjustments_created_at ON public.stock_adjustments USING btree (created_at);


--
-- Name: idx_stock_adjustments_item_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_stock_adjustments_item_id ON public.stock_adjustments USING btree (item_id);


--
-- Name: idx_stock_adjustments_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_stock_adjustments_user_id ON public.stock_adjustments USING btree (user_id);


--
-- Name: items items_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: items items_rack_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_rack_id_fkey FOREIGN KEY (rack_id) REFERENCES public.racks(id);


--
-- Name: racks racks_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT racks_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);


--
-- Name: sale_items sale_items_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- Name: sale_items sale_items_sale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_sale_id_fkey FOREIGN KEY (sale_id) REFERENCES public.sales(id) ON DELETE CASCADE;


--
-- Name: sales sales_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: stock_adjustments stock_adjustments_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_adjustments
    ADD CONSTRAINT stock_adjustments_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.items(id);


--
-- Name: stock_adjustments stock_adjustments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_adjustments
    ADD CONSTRAINT stock_adjustments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: users users_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- PostgreSQL database dump complete
--

\unrestrict WzY8voleUmHuD7wmDVSih8kfCIjUbeh3UhFHvfbfN0XZDusWkfPtGy31WBv69X4

