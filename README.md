Code4Brownies is a software package that supports active learning in programming and algorithms courses.  It allows instructors and students to share codes instantaneously during a studying session.  Features include:

- Teacher can copy a "scaffold" or partially-complete code to students' virtual whiteboards for them to work on.
- Teacher can provide additional hints for each scaffold.
- Teacher can provide scaffolds to a targeted group of students.
- Students work on the scaffolds and then share their code.
- Teacher can give virtual brownie points to students' shared work.
- Teacher can administer quizzes.
- Teacher can administer polls.
- Students can check in. Attendance can be done quickly.

Code4Brownies is designed to support several pedagogical strategies and practices such as scaffolding, guided instruction, differentiated instruction, early-and-often assessment, and early intervention.

## Architectural diagram

Teacher-students communication is achieved via a server running on the teacher's laptop that talks to Sublime Text IDEs installed on students' and teacher's laptops.  At the beginning of a lecture, the teacher starts the server.  When the lecture is over, the teacher stops the server.  The software is contained within the teacher's machine and the active-learning experience is contained within each studying session.

<img src="diagram.png" width=70% align="middle">

## Student Guide

[Click here](STUDENT.md)

## Teacher Guide

[Click here](TEACHER.md)





