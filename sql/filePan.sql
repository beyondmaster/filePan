--USE [test]
--GO
/****** Object:  Table [dbo].[file]    Script Date: 2018/3/16 下午5:45:32 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[file](
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[Name] [varchar](256) NULL,
	[Type] [varchar](16) NULL,
	[Path] [varchar](64) NULL,
	[Size] [decimal](18, 0) NULL,
	[CreatedBy] [int] NULL,
	[CreatedTime] [datetime] NULL,
	[PartitionId] [int] NULL,
	[FolderId] [int] NULL,
 CONSTRAINT [PK__folder_c__C2FABF93879D972D] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[folder]    Script Date: 2018/3/16 下午5:45:33 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[folder](
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[Name] [varchar](32) NULL,
	[Path] [varchar](256) NULL,
	[CreatedBy] [int] NULL,
	[CreatedTime] [datetime] NULL,
	[PartitionId] [int] NULL,
	[FolderId] [int] NULL,
 CONSTRAINT [PK__partitio__99E551619CC604EC] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[partition]    Script Date: 2018/3/16 下午5:45:33 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[partition](
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[Name] [varchar](32) NOT NULL,
	[CreatedBy] [int] NULL,
	[CreatedTime] [smalldatetime] NULL,
	[IsTop] [tinyint] NULL,
 CONSTRAINT [PK__partitio__99E5516100EA2E51] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[file] ADD  CONSTRAINT [DF_file_CreatedBy]  DEFAULT ((0)) FOR [CreatedBy]
GO
ALTER TABLE [dbo].[file] ADD  CONSTRAINT [DF_file_CreatedTime]  DEFAULT (getdate()) FOR [CreatedTime]
GO
ALTER TABLE [dbo].[folder] ADD  CONSTRAINT [DF_folder_CreatedBy]  DEFAULT ((0)) FOR [CreatedBy]
GO
ALTER TABLE [dbo].[folder] ADD  CONSTRAINT [DF_folder_CreatedTime]  DEFAULT (getdate()) FOR [CreatedTime]
GO
ALTER TABLE [dbo].[partition] ADD  CONSTRAINT [DF_partition_CreatedBy]  DEFAULT ((0)) FOR [CreatedBy]
GO
ALTER TABLE [dbo].[partition] ADD  CONSTRAINT [DF_partition_CreatedTime]  DEFAULT (getdate()) FOR [CreatedTime]
GO
ALTER TABLE [dbo].[partition] ADD  CONSTRAINT [DF_partition_IsTop]  DEFAULT ((0)) FOR [IsTop]
GO
