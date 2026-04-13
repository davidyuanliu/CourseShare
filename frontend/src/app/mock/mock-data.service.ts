export interface Course {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  name: string;
}

export interface Note {
  ID: number;
  CreatedAt: string;
  title: string;
  content: string;
  courseId: number;
  userId: number;
  author: {
    id: number;
    email: string;
  };
}
