// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package skills

import "strings"

// DataEngineerSkill implements the native skill for data engineer role.
type DataEngineerSkill struct {
	workspace string
}

// NewDataEngineerSkill creates a new DataEngineerSkill instance.
func NewDataEngineerSkill(workspace string) *DataEngineerSkill {
	return &DataEngineerSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (d *DataEngineerSkill) Name() string {
	return "data_engineer"
}

// Description returns a brief description of the skill.
func (d *DataEngineerSkill) Description() string {
	return "Data engineering expert: ETL pipelines, data warehouses, streaming, data quality."
}

// GetInstructions returns the complete data engineering protocol for the LLM.
func (d *DataEngineerSkill) GetInstructions() string {
	return dataEngineerInstructions
}

// GetAntiPatterns returns common data engineering anti-patterns to avoid.
func (d *DataEngineerSkill) GetAntiPatterns() string {
	return dataEngineerAntiPatterns
}

// GetExamples returns concrete data engineering examples.
func (d *DataEngineerSkill) GetExamples() string {
	return dataEngineerExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (d *DataEngineerSkill) BuildSkillContext() string {
	parts := make([]string, 0, 11)

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "📊 NATIVE SKILL: Data Engineer")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(
		parts,
		"**ROLE:** Expert Data Engineer specializing in reliable data pipelines, warehouses, and streaming systems.",
	)
	parts = append(parts, "")
	parts = append(parts, d.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, d.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, d.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (d *DataEngineerSkill) BuildSummary() string {
	return `<skill name="data_engineer" type="native">
  <purpose>Data engineering expert — ETL, data warehouse, streaming, data quality</purpose>
  <pattern>Use for data pipelines, data modeling, ETL, data quality, warehousing</pattern>
  <stacks>Spark, dbt, Airflow, Kafka, Snowflake, BigQuery, Redshift</stacks>
  <practices>ELT, Data Mesh, Data Quality, Schema Evolution</practices>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const dataEngineerInstructions = `## CORE RESPONSIBILITIES

### 1. Data Pipeline Development
- Design ETL/ELT pipelines
- Implement data transformations
- Schedule pipeline execution
- Monitor pipeline health
- Handle failures gracefully

### 2. Data Warehouse Design
- Model data for analytics (star schema)
- Implement slowly changing dimensions
- Optimize query performance
- Manage data partitions
- Design data marts

### 3. Streaming Data
- Implement real-time pipelines
- Process event streams
- Handle late-arriving data
- Manage state in streaming
- Ensure exactly-once processing

### 4. Data Quality
- Implement data validation
- Monitor data quality metrics
- Handle data anomalies
- Document data lineage
- Ensure data governance

### 5. Data Infrastructure
- Manage data storage (S3, GCS)
- Configure data processing (Spark, Flink)
- Optimize data formats (Parquet, Avro)
- Implement data cataloging
- Ensure data security

## TECHNOLOGY STACK

### Data Processing
- **Batch**: Apache Spark, dbt, Pandas
- **Streaming**: Apache Kafka, Flink, Spark Streaming

### Data Storage
- **Warehouse**: Snowflake, BigQuery, Redshift
- **Lake**: S3, GCS, ADLS
- **Database**: PostgreSQL, MySQL, MongoDB

### Orchestration
- Airflow, Prefect, Dagster, Luigi

### Data Quality
- Great Expectations, dbt tests, Soda

## BEST PRACTICES

### Data Pipeline Design
1. Extract: Pull from source systems
2. Load: Store in raw format
3. Transform: Clean and model data
4. Serve: Expose to consumers

### Data Modeling
- Star schema for analytics
- Normalize for OLTP
- Denormalize for performance
- Document data dictionary

### Data Quality Checks
- Completeness (no missing data)
- Accuracy (matches source)
- Consistency (across systems)
- Timeliness (fresh data)
- Validity (within expected ranges)

## QUALITY CHECKLIST

Before considering a pipeline complete:

- [ ] Data quality tests passing
- [ ] Monitoring alerts configured
- [ ] Documentation updated
- [ ] Backfill procedure defined
- [ ] Recovery procedure tested
- [ ] Performance benchmarks met
- [ ] Security controls in place
- [ ] Cost estimate calculated

## DATA GOVERNANCE

### Data Classification
- Public: Can be shared openly
- Internal: For internal use only
- Confidential: Sensitive business data
- PII: Personal identifiable information

### Access Control
- Role-based access (RBAC)
- Least privilege principle
- Audit logging enabled
- Regular access reviews
`

const dataEngineerAntiPatterns = `## DATA ENGINEERING ANTI-PATTERNS

### ❌ No Data Validation
` + bt + bt + bt + `python
# BAD - Trust all input data
def process_data(data):
    return transform(data)

# GOOD - Validate before processing
def process_data(data):
    assert 'id' in data, "Missing required field: id"
    assert isinstance(data['value'], (int, float)), "Value must be numeric"
    assert data['value'] >= 0, "Value must be non-negative"
    return transform(data)
` + bt + bt + bt + `

### ❌ Silent Failures
` + bt + bt + bt + `python
# BAD - Swallow errors
try:
    process_batch(data)
except Exception:
    pass  # Pipeline continues with missing data!

# GOOD - Handle errors explicitly
try:
    process_batch(data)
except Exception as e:
    logger.error(f"Batch processing failed: {e}")
    send_alert("Pipeline failure", str(e))
    raise  # Fail the pipeline
` + bt + bt + bt + `

### ❌ No Monitoring
` + bt + bt + bt + `
BAD:  Pipeline runs nightly
      No alerts on failure
      Data issues discovered by users

GOOD: Pipeline runs nightly
      Alerts on failure within 5 minutes
      Data quality dashboard
      Proactive issue detection
` + bt + bt + bt + `

### ❌ Hardcoded Paths/Credentials
` + bt + bt + bt + `python
# BAD
S3_PATH = "s3://my-bucket/data/"
DB_PASSWORD = "secret123"

# GOOD
S3_PATH = os.environ["DATA_S3_PATH"]
DB_PASSWORD = os.environ["DB_PASSWORD"]
` + bt + bt + bt + `

### ❌ No Data Lineage
` + bt + bt + bt + `
BAD:  "Where did this number come from?"
      "Which pipeline created this table?"
      "What transformations were applied?"

GOOD: Data lineage documented
      Source systems tracked
      Transformations logged
      Impact analysis possible
` + bt + bt + bt + `

### ❌ Ignoring Data Skew
` + bt + bt + bt + `python
# BAD - One partition has 90% of data
# Join takes 4 hours due to skew

# GOOD - Salt keys to distribute
df.withColumn("salt", F.rand() * 10)
  .groupBy("salt", "key")
  .agg(...)
` + bt + bt + bt + `

### ❌ Processing PII Without Encryption
` + bt + bt + bt + `python
# BAD - PII in plain text
df.write.parquet("s3://data/users/")

# GOOD - Encrypt PII
df.withColumn("email", encrypt(col("email"))) \
  .write \
  .parquet("s3://data/users/")
` + bt + bt + bt + `

### ❌ No Backfill Strategy
` + bt + bt + bt + `
BAD:  Pipeline starts today
      Historical data missing
      "We can't analyze trends"

GOOD: Backfill procedure defined
      Historical data loaded
      Incremental updates work
      Can reprocess if needed
` + bt + bt + bt + `
`

const dataEngineerExamples = `## EXAMPLE 1: CREATE DBT MODEL

**Request:** "Create a dbt model for daily user metrics"

**Expert Response:**

` + bt + bt + bt + `sql
-- models/marts/user_daily_metrics.sql
{{
  config(
    materialized='table',
    tags=['users', 'daily']
  )
}}

with user_events as (
    select * from {{ ref('stg_user_events') }}
),

daily_aggregates as (
    select
        date_trunc('day', event_timestamp) as event_date,
        user_id,
        count(*) as total_events,
        count(distinct session_id) as total_sessions,
        sum(case when event_type = 'page_view' then 1 else 0 end) as page_views,
        sum(case when event_type = 'click' then 1 else 0 end) as clicks,
        sum(case when event_type = 'purchase' then 1 else 0 end) as purchases,
        sum(purchase_amount) as total_revenue
    from user_events
    group by 1, 2
),

user_demographics as (
    select * from {{ ref('dim_users') }}
)

select
    da.event_date,
    da.user_id,
    u.country,
    u.signup_date,
    u.subscription_tier,
    da.total_events,
    da.total_sessions,
    da.page_views,
    da.clicks,
    da.purchases,
    da.total_revenue,
    case
        when da.total_events > 10 then 'highly_active'
        when da.total_events > 3 then 'active'
        else 'inactive'
    end as activity_level
from daily_aggregates da
left join user_demographics u on da.user_id = u.user_id
` + bt + bt + bt + `

## EXAMPLE 2: CREATE AIRFLOW DAG

**Request:** "Create an Airflow DAG for ETL pipeline"

**Expert Response:**

` + bt + bt + bt + `python
# dags/etl_pipeline.py
from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.providers.amazon.aws.operators.s3 import S3Operator
from datetime import datetime, timedelta

default_args = {
    'owner': 'data-team',
    'depends_on_past': False,
    'email_on_failure': True,
    'email_on_retry': False,
    'retries': 3,
    'retry_delay': timedelta(minutes=5),
}

def extract_data(**context):
    """Extract data from source API"""
    import requests

    execution_date = context['execution_date']
    response = requests.get(
        f'https://api.example.com/data?date={execution_date}'
    )
    response.raise_for_status()

    # Save to S3
    with open('/tmp/data.json', 'w') as f:
        f.write(response.text)

    return '/tmp/data.json'

def transform_data(**context):
    """Transform and clean data"""
    import pandas as pd

    ti = context['ti']
    input_path = ti.xcom_pull(task_ids='extract')

    df = pd.read_json(input_path)

    # Clean data
    df = df.dropna(subset=['id', 'value'])
    df = df[df['value'] > 0]

    # Save transformed
    output_path = '/tmp/transformed.parquet'
    df.to_parquet(output_path)

    return output_path

def load_data(**context):
    """Load data to warehouse"""
    import pyarrow.parquet as pq

    ti = context['ti']
    input_path = ti.xcom_pull(task_ids='transform')

    table = pq.read_table(input_path)

    # Load to Snowflake/BigQuery/etc.
    # ... loading logic ...

    return f"Loaded {table.num_rows} rows"

with DAG(
    'etl_pipeline',
    default_args=default_args,
    description='Daily ETL pipeline',
    schedule_interval='@daily',
    start_date=datetime(2026, 1, 1),
    catchup=True,
    tags=['etl', 'daily'],
) as dag:

    extract = PythonOperator(
        task_id='extract',
        python_callable=extract_data,
    )

    transform = PythonOperator(
        task_id='transform',
        python_callable=transform_data,
    )

    load = PythonOperator(
        task_id='load',
        python_callable=load_data,
    )

    extract >> transform >> load
` + bt + bt + bt + `
`
