CREATE TABLE segments (
                        title varchar(100) PRIMARY KEY NOT NULL,
                        description text DEFAULT '-',
                        created_at TIMESTAMP WITH TIME ZONE NOT NULL,
                        updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE users_segments (
                        id BIGSERIAL PRIMARY KEY NOT NULL,
                        user_id int NOT NULL,
                        seg_title varchar(100) NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE NOT NULL,
                        updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);