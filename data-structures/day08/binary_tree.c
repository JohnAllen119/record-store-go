#include <stdio.h>
#include <stdlib.h>
typedef struct Node
{
    int data;
    struct Node *left;
    struct Node *right;

} Node;

Node *createNode(int value)
{
    Node *newNode = malloc(sizeof(Node));
    if (newNode == NULL)
    {
        fprintf(stderr, "memory allocation failed\n");
        exit(EXIT_FAILURE);
    }
    newNode->data = value;
    newNode->left = NULL;
    newNode->right = NULL;
    return newNode;
}
void freeTree(Node *root)
{
    if (root == NULL)
        return;

    freeTree(root->left);
    freeTree(root->right);
    free(root);
}
void preorder(Node *root)
{
    if (root == NULL)
        return;
    printf("%d\n", root->data);
    preorder(root->left);
    preorder(root->right);
}
void inorder(Node *root)
{
    if (root == NULL)
        return;
    inorder(root->left);
    printf("%d\n", root->data);
    inorder(root->right);
}
void postorder(Node *root)
{
    if (root == NULL)
        return;
    postorder(root->left);
    postorder(root->right);
    printf("%d\n", root->data);
}

void levelOrder(Node *root) {
    if (root == NULL) return;

    Node *queue[100];
    int front = 0;
    int rear = 0;

    queue[rear] = root;
    rear++;
    
    while(front<rear){
        Node* current=queue[front];
        front+=1;
        printf("%d\n",current->data);
        if(current->left){
            queue[rear]=current->left;
            rear++;
        }
        if(current->right){
            queue[rear]=current->right;
            rear++;
        }
    }
}
int main(void)
{
    Node *root = createNode(10);
    Node *lroot = createNode(20);
    Node *rroot = createNode(30);
    root->left = lroot;
    root->right = rroot;
    root->left->left = createNode(40);
    root->left->right = createNode(50);
    levelOrder(root);
    freeTree(root);

    return 0;
}