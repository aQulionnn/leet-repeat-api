```mermaid
erDiagram

progress {
    int id pk
    enum perceived_difficulty
    datetime last_solved_at_utc 
    datetime next_review_at_utc 
    enum status 
    int problem_question_id 
    string problem_question 
    enum problem_difficulty 
    string problem_list_name 
}
```