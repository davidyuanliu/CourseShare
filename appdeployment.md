# Application Deployment Guide

## Prerequisites
- Node.js
- Java (JDK 17+)
- MongoDB

## Steps

### Backend
1. Navigate to backend folder
2. Run:
   mvn spring-boot:run

### Frontend
1. Navigate to frontend folder
2. Run:
   npm install
   ng serve

### Access
Frontend: http://localhost:4200  
Backend: http://localhost:8080  

## Deployment (Optional)
- Backend: Deploy using AWS EC2 or Render
- Frontend: Deploy using Netlify or Vercel