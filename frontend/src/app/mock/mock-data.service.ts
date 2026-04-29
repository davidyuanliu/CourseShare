export interface Course {
  id: number;
  CreatedAt: string;
  UpdatedAt: string;
  name: string;
}

export interface Note {
  id: number;
  CreatedAt: string;
  UpdatedAt?: string;
  title: string;
  content: string;
  courseId: number;
  course?: Course;
  userId: number;
  author?: {
    id: number;
    email: string;
  };
  authorName?: string;
  tags?: string[];
  helpfulCount?: number;
  isHelpful?: boolean;
  isSaved?: boolean;
}
