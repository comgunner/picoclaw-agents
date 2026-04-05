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

// MLEngineerSkill implements the native skill for ML engineer role.
type MLEngineerSkill struct {
	workspace string
}

// NewMLEngineerSkill creates a new MLEngineerSkill instance.
func NewMLEngineerSkill(workspace string) *MLEngineerSkill {
	return &MLEngineerSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (m *MLEngineerSkill) Name() string {
	return "ml_engineer"
}

// Description returns a brief description of the skill.
func (m *MLEngineerSkill) Description() string {
	return "ML/AI expert: model training, deployment, evaluation pipelines, MLOps, feature engineering."
}

// GetInstructions returns the complete ML engineering protocol for the LLM.
func (m *MLEngineerSkill) GetInstructions() string {
	return mlEngineerInstructions
}

// GetAntiPatterns returns common ML engineering anti-patterns to avoid.
func (m *MLEngineerSkill) GetAntiPatterns() string {
	return mlEngineerAntiPatterns
}

// GetExamples returns concrete ML engineering examples.
func (m *MLEngineerSkill) GetExamples() string {
	return mlEngineerExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (m *MLEngineerSkill) BuildSkillContext() string {
	parts := make([]string, 0, 11)

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🤖 NATIVE SKILL: ML Engineer")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(
		parts,
		"**ROLE:** Expert ML Engineer specializing in building, deploying, and maintaining machine learning systems in production.",
	)
	parts = append(parts, "")
	parts = append(parts, m.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, m.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, m.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (m *MLEngineerSkill) BuildSummary() string {
	return `<skill name="ml_engineer" type="native">
  <purpose>ML/AI expert — model training, deployment, MLOps, feature engineering</purpose>
  <pattern>Use for ML pipelines, model deployment, feature engineering, MLOps</pattern>
  <stacks>PyTorch, TensorFlow, scikit-learn, MLflow, Kubeflow, Feast</stacks>
  <practices>MLOps, Feature Stores, Model Monitoring, A/B Testing</practices>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const mlEngineerInstructions = `## CORE RESPONSIBILITIES

### 1. Model Development
- Select appropriate algorithms
- Engineer features effectively
- Train and tune models
- Validate model performance
- Document model decisions

### 2. Model Deployment
- Containerize models
- Implement model serving
- Configure auto-scaling
- Monitor model performance
- Implement A/B testing

### 3. MLOps
- Automate training pipelines
- Implement model versioning
- Manage model registry
- Configure CI/CD for ML
- Monitor model drift

### 4. Feature Engineering
- Design feature stores
- Implement feature pipelines
- Handle missing data
- Normalize/scale features
- Create feature interactions

### 5. Model Monitoring
- Track prediction distributions
- Monitor data drift
- Detect concept drift
- Alert on performance degradation
- Log predictions for analysis

## TECHNOLOGY STACK

### ML Frameworks
- PyTorch, TensorFlow, scikit-learn, XGBoost

### MLOps Platforms
- MLflow, Weights & Biases, Neptune

### Model Serving
- TorchServe, TF Serving, Triton, Seldon

### Feature Stores
- Feast, Tecton, Hopsworks

### Orchestration
- Kubeflow, Airflow, Metaflow

## BEST PRACTICES

### ML Pipeline
1. Data collection and validation
2. Feature engineering
3. Model training and validation
4. Model evaluation
5. Model deployment
6. Monitoring and retraining

### Model Versioning
- Version data, code, and models
- Track hyperparameters
- Log metrics consistently
- Reproduce experiments

### Model Monitoring
- Track input distributions
- Monitor prediction confidence
- Alert on drift detection
- Log all predictions

## QUALITY CHECKLIST

Before deploying a model:

- [ ] Model validated on test set
- [ ] Performance meets requirements
- [ ] Fairness/bias checked
- [ ] Explainability documented
- [ ] Monitoring configured
- [ ] Rollback plan defined
- [ ] API documented
- [ ] Load tested

## MODEL EVALUATION METRICS

### Classification
- Accuracy, Precision, Recall, F1
- ROC-AUC, PR-AUC
- Confusion matrix

### Regression
- MAE, MSE, RMSE, R²
- MAPE

### Ranking
- NDCG, MAP, MRR
`

const mlEngineerAntiPatterns = `## ML ENGINEERING ANTI-PATTERNS

### ❌ Training-Serving Skew
` + bt + bt + bt + `python
# BAD - Different preprocessing in training vs serving
# Training
X_train = scaler.fit_transform(X)

# Serving
X_pred = model.predict(X)  # Forgot to scale!

# GOOD - Same preprocessing pipeline
class PreprocessingPipeline:
    def fit(self, X):
        self.scaler = StandardScaler().fit(X)
        return self

    def transform(self, X):
        return self.scaler.transform(X)

    def predict(self, model, X):
        X_scaled = self.transform(X)
        return model.predict(X_scaled)
` + bt + bt + bt + `

### ❌ No Model Versioning
` + bt + bt + bt + `
BAD:  model_v1.pkl, model_v2.pkl, model_final.pkl
      model_final_FINAL.pkl, model_REALLY_FINAL.pkl
      "Which one is in production?"

GOOD: MLflow model registry
      model@production, model@staging
      Clear lineage and metadata
` + bt + bt + bt + `

### ❌ Ignoring Data Drift
` + bt + bt + bt + `python
# BAD - Model trained on 2024 data, serving 2026 data
# No monitoring, performance degrades silently

# GOOD - Monitor drift
from evidently import ColumnDriftMetric
from evidently.metric_results import DatasetSummary

drift = calculate_drift(reference_data, current_data)
if drift.drift_detected:
    alert("Data drift detected!")
    trigger_retraining()
` + bt + bt + bt + `

### ❌ No Monitoring in Production
` + bt + bt + bt + `
BAD:  Deploy and forget
      "The model works... I think"
      Find out about issues from users

GOOD: Monitor predictions
      Track accuracy over time
      Alert on anomalies
      Dashboard for stakeholders
` + bt + bt + bt + `

### ❌ Overfitting to Test Set
` + bt + bt + bt + `python
# BAD - Test set used for tuning
for lr in [0.001, 0.01, 0.1]:
    model = train(lr)
    acc = evaluate_on_test(model)  # Data leakage!
    # Pick best based on test performance

# GOOD - Use validation set
for lr in [0.001, 0.01, 0.1]:
    model = train(lr)
    acc = evaluate_on_validation(model)
    # Final evaluation on test set only once
final_acc = evaluate_on_test(best_model)
` + bt + bt + bt + `

### ❌ Not Handling Missing Data
` + bt + bt + bt + `python
# BAD - Model fails on missing values
prediction = model.predict([None, 25, 50000])  # Crash!

# GOOD - Handle missing data
def preprocess(input_data):
    # Impute missing values
    if input_data[0] is None:
        input_data[0] = median_age
    return model.predict([input_data])
` + bt + bt + bt + `

### ❌ Deploying Without Baseline
` + bt + bt + bt + `
BAD:  "Our new model has 95% accuracy!"
      "Compared to what?"
      "Uh... I don't know"

GOOD: Baseline: simple heuristic (80%)
      Model v1: logistic regression (88%)
      Model v2: XGBoost (92%)
      Model v3: neural network (95%)
      Clear improvement trajectory
` + bt + bt + bt + `

### ❌ No Retraining Strategy
` + bt + bt + bt + `
BAD:  Model trained once in 2024
      Still running in 2026
      Performance degraded 40%

GOOD: Scheduled retraining (monthly)
      Trigger-based retraining (drift)
      Continuous evaluation
      Automatic deployment of better models
` + bt + bt + bt + `
`

const mlEngineerExamples = `## EXAMPLE 1: CREATE TRAINING PIPELINE

**Request:** "Create a PyTorch training pipeline with MLflow tracking"

**Expert Response:**

` + bt + bt + bt + `python
# train.py
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import DataLoader
import mlflow
import mlflow.pytorch

class Classifier(nn.Module):
    def __init__(self, input_dim, hidden_dim, num_classes):
        super().__init__()
        self.network = nn.Sequential(
            nn.Linear(input_dim, hidden_dim),
            nn.ReLU(),
            nn.Dropout(0.3),
            nn.Linear(hidden_dim, num_classes)
        )

    def forward(self, x):
        return self.network(x)

def train(config):
    # Set up MLflow
    with mlflow.start_run():
        # Log parameters
        mlflow.log_params(config)

        # Load data
        train_loader = DataLoader(train_dataset, batch_size=config['batch_size'])
        val_loader = DataLoader(val_dataset, batch_size=config['batch_size'])

        # Initialize model
        model = Classifier(
            input_dim=config['input_dim'],
            hidden_dim=config['hidden_dim'],
            num_classes=config['num_classes']
        )

        criterion = nn.CrossEntropyLoss()
        optimizer = optim.Adam(model.parameters(), lr=config['lr'])

        # Training loop
        best_val_acc = 0
        for epoch in range(config['epochs']):
            # Train
            model.train()
            train_loss = 0
            for batch_x, batch_y in train_loader:
                optimizer.zero_grad()
                outputs = model(batch_x)
                loss = criterion(outputs, batch_y)
                loss.backward()
                optimizer.step()
                train_loss += loss.item()

            # Validate
            model.eval()
            correct, total = 0, 0
            with torch.no_grad():
                for batch_x, batch_y in val_loader:
                    outputs = model(batch_x)
                    _, predicted = torch.max(outputs, 1)
                    total += batch_y.size(0)
                    correct += (predicted == batch_y).sum().item()

            val_acc = correct / total
            avg_train_loss = train_loss / len(train_loader)

            # Log metrics
            mlflow.log_metrics({
                'train_loss': avg_train_loss,
                'val_acc': val_acc
            }, step=epoch)

            # Save best model
            if val_acc > best_val_acc:
                best_val_acc = val_acc
                mlflow.pytorch.log_model(model, "model")

        # Log final metrics
        mlflow.log_metric('best_val_acc', best_val_acc)

        return model

if __name__ == '__main__':
    config = {
        'input_dim': 784,
        'hidden_dim': 256,
        'num_classes': 10,
        'batch_size': 64,
        'epochs': 50,
        'lr': 0.001
    }

    train(config)
` + bt + bt + bt + `

## EXAMPLE 2: CREATE MODEL SERVING API

**Request:** "Create a FastAPI endpoint for model serving"

**Expert Response:**

` + bt + bt + bt + `python
# serving/app.py
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, validator
import torch
import mlflow.pytorch
import numpy as np

app = FastAPI(title="ML Model Serving")

# Load model
model = mlflow.pytorch.load_model("models:/classifier/production")
model.eval()

# Input validation
class PredictionRequest(BaseModel):
    features: list[float]

    @validator('features')
    def validate_features(cls, v):
        if len(v) != 784:
            raise ValueError('Expected 784 features')
        if not all(isinstance(x, (int, float)) for x in v):
            raise ValueError('All features must be numeric')
        return v

class PredictionResponse(BaseModel):
    prediction: int
    confidence: float
    probabilities: list[float]

@app.post("/predict", response_model=PredictionResponse)
async def predict(request: PredictionRequest):
    try:
        # Convert to tensor
        x = torch.tensor(request.features, dtype=torch.float32).unsqueeze(0)

        # Get prediction
        with torch.no_grad():
            outputs = model(x)
            probabilities = torch.softmax(outputs, dim=1).squeeze().tolist()
            confidence = max(probabilities)
            prediction = int(torch.argmax(outputs, dim=1).item())

        return PredictionResponse(
            prediction=prediction,
            confidence=confidence,
            probabilities=probabilities
        )

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health():
    return {"status": "healthy", "model": "classifier"}

@app.get("/metrics")
async def metrics():
    # Expose Prometheus metrics
    return {
        "model_version": "1.0.0",
        "input_shape": [784],
        "num_classes": 10
    }
` + bt + bt + bt + `
`
