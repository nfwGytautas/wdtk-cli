package templates

// PUBLIC TYPES
// ========================================================================

/*
K8S template data
*/
type K8STemplateData struct {
	ProjectName string
	Service     string
}

/*
Template for deployment file
*/
const K8STemplate = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ProjectName}}-{{.Service}}-service
spec:
  selector:
    matchLabels:
      app: {{.ProjectName}}-{{.Service}}-service
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: {{.ProjectName}}-{{.Service}}-service
    spec:
      containers:
      - image: {{.ProjectName}}/{{.Service}}-service:0.0.0
        name: {{.ProjectName}}-{{.Service}}-service
        imagePullPolicy: Never
        resources:
          limits:
            memory: "500M"
            cpu: "50m"
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: API_SECRET
          valueFrom:
            secretKeyRef:
              name: mstk-project-secret
              key: Secret
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ProjectName}}-{{.Service}}-balancer
spec:
  selector:
    matchLabels:
      app: {{.ProjectName}}-{{.Service}}-balancer
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: {{.ProjectName}}-{{.Service}}-balancer
    spec:
      containers:
      - image: {{.ProjectName}}/{{.Service}}-balancer:0.0.0
        name: {{.ProjectName}}-{{.Service}}-balancer
        imagePullPolicy: Never
        resources:
          limits:
            memory: "500M"
            cpu: "50m"
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: API_SECRET
          valueFrom:
            secretKeyRef:
              name: mstk-project-secret
              key: Secret
---
apiVersion: v1
kind: Service
metadata:
  name: {{.ProjectName}}-{{.Service}}-service
spec:
  ports:
    - port: 8080
      targetPort: 8080
  selector:
    app: {{.ProjectName}}-{{.Service}}-service
---
apiVersion: v1
kind: Service
metadata:
  name: {{.ProjectName}}-{{.Service}}-balancer
spec:
  ports:
    - port: 8080
      targetPort: 8080
  selector:
    app: {{.ProjectName}}-{{.Service}}-balancer
`
